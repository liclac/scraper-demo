package lib

import (
	"golang.org/x/net/html"
)

// Gets an attribute from an HTML node, or "" if it doesn't exist.
func attr(node *html.Node, name string) string {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}
