package rssnotify

import (
	"fmt"
	"os"
	"strings"
)

func Run() error {
	notifyURL := os.Getenv("NOTIFY_URL")
	if notifyURL == "" {
		return fmt.Errorf("NOTIFY_URL is not set")
	}
	feedURLs := strings.Split(os.Getenv("FEED_URLS"), ";")
	if len(feedURLs) == 0 {
		return fmt.Errorf("FEED_URLS is not set")
	}

	reader := newReadClient()
	postsCh := reader.readPostsWithinAWeek(feedURLs)
	msg := makeMsg(postsCh)

	if err := notify(notifyURL, msg); err != nil {
		return err
	}
	return nil
}

type post struct {
	title string
	link  string
}
