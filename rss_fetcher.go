package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// RSS represents the structure of an RSS feed
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the main content of the RSS feed
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       []Item `xml:"item"`
}

// Item represents a single entry in the RSS feed
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type FeedWorker struct {
	interval            time.Duration
	batchSize           int
	getNextFeedsToFetch func(int32) ([]Feed, error) // 用 int32 是因為 sqlc 產生的 function 用 int32
}

func fetchRSS(url string) (RSS, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return RSS{}, fmt.Errorf("error fetching RSS feed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSS{}, fmt.Errorf("error reading response body: %v", err)
	}

	rss := RSS{}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return RSS{}, fmt.Errorf("error parsing XML: %v", err)
	}

	fmt.Printf("Title: %s\n", rss.Channel.Title)
	fmt.Printf("Description: %s\n", rss.Channel.Description)
	fmt.Printf("Number of items: %d\n", len(rss.Channel.Items))

	// Print the first item as an example
	if len(rss.Channel.Items) > 0 {
		fmt.Printf("First item title: %s\n", rss.Channel.Items[0].Title)
	}

	fmt.Println("-------------------------")
	return rss, nil
}

func NewFeedWorker(
	interval time.Duration,
	batchSize int,
	getNextFeedsToFetch func(int32) ([]Feed, error),
) *FeedWorker {
	return &FeedWorker{
		interval:  interval,
		batchSize: batchSize,
		getNextFeedsToFetch: func(batchSize int32) ([]Feed, error) {
			return getNextFeedsToFetch(batchSize)
		},
	}
}

func (w *FeedWorker) Start() {
	log.Println("Starting feed worker")
	w.processBatch()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for range ticker.C {
		w.processBatch()
	}
}

func (w *FeedWorker) processBatch() {
	log.Printf("Processing batch of up to %d feeds\n", w.batchSize)
	feeds, err := w.getNextFeedsToFetch(int32(w.batchSize))
	if err != nil {
		log.Println("processBatch", err)
		return
	}

	var wg sync.WaitGroup
	for _, v := range feeds {
		wg.Add(1)
		go func(f Feed) {
			defer wg.Done()
			w.processFeed(f)
		}(v)
	}
	wg.Wait()
	log.Println("Finished processing batch")
}

func (cfg *apiConfig) getNextFeedsToFetch(batchSize int32) ([]Feed, error) {
	ctx := context.Background()
	dbFeeds, err := cfg.DB.GetNextFeedsToFetch(ctx, batchSize)
	if err != nil {
		fmt.Println("getNextFeeds", err)
		return nil, err
	}

	feeds := make([]Feed, len(dbFeeds))
	for i, v := range dbFeeds {
		feeds[i] = databaseFeedToFeed(v)
	}

	return feeds, nil
}

func (w *FeedWorker) processFeed(feed Feed) {
	log.Printf("Processing feed ID %s: %s\n", feed.ID, feed.URL)
	res, err := fetchRSS(feed.URL)
	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range res.Channel.Items {
		log.Printf("%s - %s\n", res.Channel.Title, item.Title)
	}
}
