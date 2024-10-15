# Builder stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY src .
RUN apk add --no-cache upx && \
    go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate -path ./frontend/templates/ && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o blogo . && \
    upx blogo && \
    chmod a+rx blogo && \
    mkdir -p /app/articles

# Node stage
FROM node:alpine AS node
RUN apk update && apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY src/frontend/tailwind.config.js .
COPY src/frontend/package.json src/frontend/package-lock.json ./
COPY src/frontend/static ./static
COPY src/frontend/templates ./templates
RUN npm ci && \
    npx tailwindcss -i ./static/css/input.css -o ./static/css/style.css --minify

FROM scratch
COPY --from=node /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Include timezone data
COPY --from=node /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /blogo
COPY --from=builder /app/blogo /blogo/blogo
COPY --from=builder /app/articles /blogo/articles
COPY ./src/frontend/static /blogo/frontend/static
COPY --from=node /app/static/css/style.css /blogo/frontend/static/css/style.css

# Ensure the app uses the right timezone by setting the TZ environment variable
ENV TZ="Europe/Warsaw"
ENV PATH="/blogo:$PATH"

EXPOSE 1337
CMD ["/blogo/blogo", "serve"]