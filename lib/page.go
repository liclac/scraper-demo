package lib

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Represents a page and any links or assets on it.
type Page struct {
	URL           *url.URL
	Assets        []Asset
	InternalLinks []*url.URL
	ExternalLinks []*url.URL
}

// Loads a URL as a Page.
func Load(ctx context.Context, u *url.URL) (*Page, error) {
	req := &http.Request{Method: http.MethodGet, URL: u}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return Parse(res.Body, req.URL)
}

// Parses HTML data into a Page structure.
func Parse(r io.Reader, self *url.URL) (*Page, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var fn func(page *Page, node *html.Node) error
	fn = func(page *Page, node *html.Node) error {
		for n := node.FirstChild; n != nil; n = n.NextSibling {
			if n.Type == html.ElementNode {
				switch n.DataAtom {
				case atom.A:
					u := page.resolveLink(attr(n, "href"))
					if u == nil {
						continue
					}
					if u.Host == page.URL.Host {
						page.InternalLinks = append(page.InternalLinks, u)
					} else {
						page.ExternalLinks = append(page.ExternalLinks, u)
					}
				case atom.Link:
					u := page.resolveLink(attr(n, "href"))
					if u == nil {
						continue
					}
					switch attr(n, "rel") {
					case "stylesheet":
						page.Assets = append(page.Assets, Asset{Type: Stylesheet, URL: u})
					}
				case atom.Script:
					src := attr(n, "src")
					if src == "" {
						continue
					}
					if u := page.resolveLink(src); u != nil {
						page.Assets = append(page.Assets, Asset{Type: Script, URL: u})
					}
				case atom.Img:
					if u := page.resolveLink(attr(n, "src")); u != nil {
						page.Assets = append(page.Assets, Asset{Type: Image, URL: u})
					}
				}
			}

			if err := fn(page, n); err != nil {
				return nil
			}
		}
		return nil
	}

	page := Page{URL: self}
	return &page, fn(&page, root)
}

// Resolves a link on a page
func (p *Page) resolveLink(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		return nil
	}
	if !u.IsAbs() {
		u = p.URL.ResolveReference(u)
	}
	u.Fragment = ""
	return u
}
