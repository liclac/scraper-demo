package lib

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPageSrc = `
<html>
<head>
	<title>Test Page</title>

	<link rel="stylesheet" href="http://example.com/somelib.css" />
	<link rel="stylesheet" href="/style.css" />
	<link rel="preload" href="/image.png" />

	<script src="http://example.com/lib.js"></script>
	<script src="/script.js"></script>
	<script>
	alert("beep boop");
	</script>
</head>
<body>
	<h1>Lorem <a name="asdf" href="#asdf">ipsum</a> dolor sit amet</h1>
	<a href="/somepage#hi"><img src="/image.png" /></a>
	<a href="http://example.com/#anchor">Hi</a>
</body>
</html>
`

func TestParse(t *testing.T) {
	u, err := url.Parse("https://lorem.ipsum/")
	if !assert.NoError(t, err) {
		return
	}

	buf := bytes.NewBufferString(testPageSrc)
	page, err := Parse(buf, u)
	if assert.NoError(t, err) {
		assert.Equal(t, "https://lorem.ipsum/", page.URL.String())

		if assert.Len(t, page.Assets, 5) {
			assert.Equal(t, Stylesheet, page.Assets[0].Type)
			assert.Equal(t, "http://example.com/somelib.css", page.Assets[0].URL.String())
			assert.Equal(t, Stylesheet, page.Assets[1].Type)
			assert.Equal(t, "https://lorem.ipsum/style.css", page.Assets[1].URL.String())
			assert.Equal(t, Script, page.Assets[2].Type)
			assert.Equal(t, "http://example.com/lib.js", page.Assets[2].URL.String())
			assert.Equal(t, Script, page.Assets[3].Type)
			assert.Equal(t, "https://lorem.ipsum/script.js", page.Assets[3].URL.String())
			assert.Equal(t, Image, page.Assets[4].Type)
			assert.Equal(t, "https://lorem.ipsum/image.png", page.Assets[4].URL.String())
		}

		if assert.Len(t, page.InternalLinks, 2) {
			assert.Equal(t, "https://lorem.ipsum/", page.InternalLinks[0].String())
			assert.Equal(t, "https://lorem.ipsum/somepage", page.InternalLinks[1].String())
		}
		if assert.Len(t, page.ExternalLinks, 1) {
			assert.Equal(t, "http://example.com/", page.ExternalLinks[0].String())
		}
	}
}
