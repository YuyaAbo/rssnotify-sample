package rssnotify

import (
	"fmt"
	"net/http"
	"strings"
)

const baseMsg = `{"text":"%s"}`

func makeMsg(ch <-chan post) string {
	var msg string
	for p := range ch {
		msg += fmt.Sprintf("%s\n%s\n\n", p.title, p.link)
	}
	if msg != "" {
		return fmt.Sprintf(baseMsg, msg)
	}
	return fmt.Sprintf(baseMsg, "no post")
}

func notify(url, msg string) error {
	body := strings.NewReader(msg)
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
