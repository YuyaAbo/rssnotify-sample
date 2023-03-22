package rssnotify

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

type readClient struct {
	*http.Client
}

func newReadClient(opts ...option) *readClient {
	cli := &readClient{http.DefaultClient}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

type option func(*readClient)

func optionHTTPClient(c *http.Client) option {
	return func(cli *readClient) {
		cli.Client = c
	}
}

func (c *readClient) readPostsWithinAWeek(feedURLs []string) <-chan post {
	var wg sync.WaitGroup
	wg.Add(len(feedURLs))

	ch := make(chan post)
	go func() {
		for _, url := range feedURLs {
			go func(url string) {
				defer wg.Done()

				parser := gofeed.NewParser() // Avoid data races
				parser.Client = c.Client
				feed, err := parser.ParseURL(url)
				if err != nil {
					log.Println(err)
					return
				}

				for _, item := range feed.Items {
					if item.PublishedParsed == nil {
						continue
					}
					now := time.Now()
					if withinAWeek(&now, item.PublishedParsed) {
						ch <- post{title: item.Title, link: item.Link}
					}
				}
			}(url)
		}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func withinAWeek(t1, t2 *time.Time) bool {
	lastWeek := t1.Add(-7 * 24 * time.Hour)
	return !t2.Before(lastWeek)
}
