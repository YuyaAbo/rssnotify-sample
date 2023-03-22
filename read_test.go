package rssnotify

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

func Test_withinAWeek(t *testing.T) {
	type args struct {
		t1 string
		t2 string
	}
	tests := map[string]struct {
		args args
		want bool
	}{
		"Just a week ago": {
			args: args{"2023/03/22 00:00", "2023/03/15 00:00"},
			want: true,
		},
		"One week and one seconds ago": {
			args: args{"2023/03/22 00:00", "2023/03/14 23:59"},
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t1 := parseTime(t, tt.args.t1)
			t2 := parseTime(t, tt.args.t2)
			if got := withinAWeek(&t1, &t2); got != tt.want {
				t.Errorf("withinAWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readClient_readPostsWithinAWeek(t *testing.T) {
	cli := newReadClient(optionHTTPClient(client(t)))
	postCh := cli.readPostsWithinAWeek([]string{"https://example.com/test1/rss", "https://example.com/test2/rss"})

	var got []post
	for p := range postCh {
		got = append(got, p)
	}

	want := []post{
		{title: "書籍『Go言語プログラミングエッセンス』を読んだ", link: "https://aboy-perry.hatenablog.com/entry/2023/03/20/001649"},
		{title: "書籍『Go言語プログラミングエッセンス』を読んだ2", link: "https://aboy-perry.hatenablog.com/entry/2023/03/20/001649"},
	}
	sort.SliceStable(got, func(i, j int) bool { return got[i].title < got[j].title })

	if !reflect.DeepEqual(got, want) {
		t.Errorf("readPostsWithinAWeek() = %v, want %v", got, want)
	}
}

func parseTime(t *testing.T, timeStr string) time.Time {
	t.Helper()
	ti, err := time.Parse("2006/01/02 15:04", timeStr)
	if err != nil {
		t.Fatal(err)
	}
	return ti
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{Transport: fn}
}

func client(t *testing.T) *http.Client {
	t.Helper()
	return newTestClient(func(req *http.Request) *http.Response {
		switch req.URL.String() {
		case "https://example.com/test1/rss":
			b := readFile(t, "testdata/feed1.rss")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(b)),
				Header:     make(http.Header),
			}
		case "https://example.com/test2/rss":
			b := readFile(t, "testdata/feed2.rss")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(b)),
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       nil,
			Header:     make(http.Header),
		}
	})
}

func readFile(t *testing.T, fileName string) []byte {
	t.Helper()
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		file.Close()
	})
	b, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
