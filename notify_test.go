package rssnotify

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func Test_makeMsg(t *testing.T) {
	noPostCh := make(chan post)
	close(noPostCh)
	if got := makeMsg(noPostCh); got != `{"text":"no post"}` {
		t.Errorf("makeMsg() = %v, want %v", got, `{"text":"no post"}`)
	}

	postsCh := make(chan post)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		go func() {
			defer wg.Done()
			postsCh <- post{title: "title1", link: "https://example.com/1"}
		}()
		go func() {
			defer wg.Done()
			postsCh <- post{title: "title2", link: "https://example.com/2"}
		}()
		go func() {
			wg.Wait()
			close(postsCh)
		}()
	}()

	// Either pattern is acceptable, as the order is not guaranteed.
	want1 := `{"text":"title1
https://example.com/1

title2
https://example.com/2

"}`
	want2 := `{"text":"title2
https://example.com/2

title1
https://example.com/1

"}`
	if got := makeMsg(postsCh); got != want1 && got != want2 {
		t.Errorf("makeMsg() = %v, want %v or %v", got, want1, want2)
	}
}

func Test_notify(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("notify() method = %v, want %v", r.Method, http.MethodPost)
		}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "test message" {
			t.Errorf("notify() send = %v, want %v", string(b), "test message")
		}
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	err := notify(ts.URL, "test message")
	if err != nil {
		t.Errorf("notify() returns error %v, want nil", err)
	}
}
