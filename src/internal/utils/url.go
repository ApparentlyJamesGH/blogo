package utils

import (
	"net/url"
	"strings"
)

func CreateURL(host string, paths ...string) string {
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}

	u, _ := url.Parse(host)
	u.Path, _ = url.JoinPath("/", paths...)

	return u.String()
}
