package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

func SearchYouTube(query string) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	searchURL := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", query)
	var videoLinks []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(searchURL),
		chromedp.Sleep(1*time.Second),
		chromedp.EvaluateAsDevTools(`
			Array.from(document.querySelectorAll("a#video-title"))
				.map(a => a.href)
				.filter(url => url.includes("watch"))
		`, &videoLinks),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape YouTube: %w", err)
	}

	return videoLinks, nil
}
