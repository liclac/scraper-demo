package lib

import (
	"net/url"
)

// An asset's type.
type AssetType int

const (
	Stylesheet AssetType = iota
	Script
	Image
)

// Represents an asset on a page.
type Asset struct {
	Type AssetType
	URL  *url.URL
}
