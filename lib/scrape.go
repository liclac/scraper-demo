package lib

import (
	"context"
	"net/url"
	"sync"
)

// Result of a Scrape call.
type ScrapeResult struct {
	Err      error
	Page     *Page
	Children []ScrapeResult
}

// Resursively scrapes a page.
func Scrape(ctx context.Context, urlStr string, depth int) (res ScrapeResult) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ScrapeResult{Err: err}
	}
	return newScraper().scrape(ctx, u, depth, 0)
}

// Internal state for a Scrape operation; we need to keep track of cyclic references.
type scraper struct {
	seen     map[string]*Page
	seenLock sync.RWMutex
}

func newScraper() *scraper {
	return &scraper{seen: make(map[string]*Page)}
}

func (s *scraper) scrape(ctx context.Context, u *url.URL, maxDepth, atDepth int) (res ScrapeResult) {
	urlStr := u.String()

	// Check the cache!
	s.seenLock.RLock()
	page := s.seen[urlStr]
	s.seenLock.RUnlock()
	if page != nil {
		res.Page = page
		return
	}

	page, err := Load(ctx, u)
	if err != nil {
		res.Page = &Page{URL: u}
		res.Err = err
		return res
	}
	res.Page = page

	// Put it in the cache (if nothing else already did)
	s.seenLock.Lock()
	if _, ok := s.seen[urlStr]; !ok {
		s.seen[urlStr] = page
	}
	s.seenLock.Unlock()

	if len(page.InternalLinks) > 0 && atDepth < maxDepth {
		children := make(chan ScrapeResult)
		for _, link := range page.InternalLinks {
			link := link
			go func() { children <- s.scrape(ctx, link, maxDepth, atDepth+1) }()
		}
		res.Children = make([]ScrapeResult, len(page.InternalLinks))
		for i := range page.InternalLinks {
			res.Children[i] = <-children
		}
		close(children)
	}
	return res
}
