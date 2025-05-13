package enddpoints_test

import (
	"fmt"
	"testing"
)

func TestLinks(t *testing.T) {
	link := NewLink("main", "All")
	link.AddChild(NewLink("none"))
	link.AddChild(NewLink("example", "example.com", "https://example.com"))

	sub := NewLink("sub", "Sub")
	sub.AddChild(NewLink("hello"))
	sub.AddChild(NewLink("login", "Login", "/login"))
	link.AddChild(sub)
	link.AddChild(NewLink("Logout"))

	fmt.Println("Links:")
	fmt.Println(link.Build(1))
}
