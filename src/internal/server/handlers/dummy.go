package handlers

import (
	"fmt"

	"github.com/kataras/iris/v12"
)

func DummyRespond(c iris.Context) {
	c.HTML(fmt.Sprintf("<h1> Hello from %s</h1>", c.Path()))
}
