package endpoints

import (
	"strings"
)

type Link struct {
	ID   string
	Name string
	Link *string

	Links []*Link
}

func NewLink(params ...string) *Link {
	if len(params) == 0 {
		panic("Expecting params {id, name, link}")
	}

	var name string
	var link *string

	if len(params) == 1 {
		name = params[0]
	} else {
		name = params[1]
	}

	if len(params) > 2 {
        link = &(params[2])
	}

	return &Link{
		ID:   params[0],
		Name: name,
		Link: link,
	}
}

func (l *Link) AddChild(c *Link) {
	if l.Links == nil {
		l.Links = make([]*Link, 0)
	}

	l.Links = append(l.Links, c)
}

func (l *Link) PrependChild(c *Link) {
	if l.Links == nil {
		l.Links = make([]*Link, 0)
	}

	l.Links = append([]*Link{c}, l.Links...)
}

func (l *Link) Build(level int) string {
	var sb strings.Builder

	sb.WriteString(strings.Repeat("-", level))
	sb.WriteRune(' ')
	sb.WriteString(l.ID)
	sb.WriteRune(',')
	sb.WriteString(l.Name)
	if l.Link != nil {
		sb.WriteRune(',')
		sb.WriteString(*l.Link)
	}
	sb.WriteString("\n")

	for _, c := range l.Links {
		sb.WriteString(c.Build(level + 1))
	}

	return sb.String()
}
