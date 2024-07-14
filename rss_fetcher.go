package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
	interval  time.Duration
	batchSize int
	cfg       *apiConfig
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
	cfg *apiConfig,
) *FeedWorker {
	return &FeedWorker{
		interval:  interval,
		batchSize: batchSize,
		cfg:       cfg,
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
	feeds, err := w.cfg.getNextFeedsToFetch(int32(w.batchSize))
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

	ctx := context.Background()
	for _, item := range res.Channel.Items {
		pubDate, err := parseRSSDateFormat(item.PubDate)
		if err != nil {
			log.Printf("cannot parse pubDate %v of item %v \n", item.PubDate, item.Title)
			continue
		}
		err = w.cfg.DB.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		})

		if err != nil {
			log.Printf("failed CreatePost %#v", item)
		}
	}
}

var rssDateFormats = []string{
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	"2006-01-02T15:04:05Z07:00", // another common variant
}

func parseRSSDateFormat(pubDate string) (time.Time, error) {
	for _, v := range rssDateFormats {
		t, err := time.Parse(v, pubDate)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unable to parse data: " + pubDate)
}
