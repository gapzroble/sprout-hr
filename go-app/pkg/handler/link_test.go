package handler_test

import (
	"fmt"
	"testing"

	"github.com/gapzroble/sprout-hr/pkg/handler"
)

func TestLinks(t *testing.T) {
	link := handler.NewLink("main", "All")
	link.AddChild(handler.NewLink("none"))
	link.AddChild(handler.NewLink("example", "example.com", "https://example.com"))

	sub := handler.NewLink("sub", "Sub")
	sub.AddChild(handler.NewLink("hello"))
	sub.AddChild(handler.NewLink("login", "Login", "/login"))
	link.AddChild(sub)
	link.AddChild(handler.NewLink("Logout"))

	fmt.Println("Links:")
	fmt.Println(link.Build(1))
}
