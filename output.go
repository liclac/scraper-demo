package main

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/liclac/scraper-demo/lib"
)

var (
	URLColor     = color.New(color.Bold)
	ErrColor     = color.New(color.FgRed)
	AssetColor   = color.New(color.FgGreen)
	IntLinkColor = color.New(color.FgBlue)
	ExtLinkColor = color.New(color.FgHiBlue)
)

func PrintResult(w io.Writer, res lib.ScrapeResult) {
	printResult(w, res, "", nil)
}

func printResult(w io.Writer, res lib.ScrapeResult, indent string, parents []*lib.Page) {
	fmt.Fprint(w, indent+URLColor.Sprintf("%s\n", res.Page.URL))
	if res.Err != nil {
		fmt.Fprint(w, indent+ErrColor.Sprintf("%s\n", res.Err))
	}
	for _, asset := range res.Page.Assets {
		fmt.Fprint(w, indent+"  "+assetEmoji(asset.Type)+AssetColor.Sprintf("\t%s\n", asset.URL))
	}
	for _, link := range res.Page.ExternalLinks {
		fmt.Fprint(w, indent+"  🔗"+ExtLinkColor.Sprintf("\t%s\n", link))
	}
	for _, link := range res.Page.InternalLinks {
		fmt.Fprint(w, indent+"  📎"+IntLinkColor.Sprintf("\t%s\n", link))
	}
	for _, child := range res.Children {
		// Don't follow links to self, that way lies madness.
		if *child.Page.URL == *res.Page.URL {
			continue
		}
		// Don't follow links to parents, it makes a mess around links to the front page.
		for _, parent := range parents {
			if *child.Page.URL == *parent.URL {
				continue
			}
		}
		fmt.Fprint(w, "\n")
		printResult(w, child, indent+"  ", append(parents, res.Page))
	}
}

func assetEmoji(t lib.AssetType) string {
	switch t {
	case lib.Stylesheet:
		return "🖌️"
	case lib.Script:
		return "🕹️"
	case lib.Image:
		return "🖼️"
	default:
		return ""
	}
}
