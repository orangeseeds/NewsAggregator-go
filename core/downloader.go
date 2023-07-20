package core

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Runner struct {
	Sources     map[string]sourceInfo
	Feeds       map[string]rssFeed
	client      *http.Client
	Failed      []string
	successLock sync.Mutex
	failedLock  sync.Mutex
	task        *sync.WaitGroup
	limit       chan bool
	cleaner     feedParser
	stop        chan bool
}

func NewRunner(concurreny int, sources map[string]sourceInfo) *Runner {
	var wg sync.WaitGroup
	return &Runner{
		Sources: sources,
		client:  &http.Client{},
		limit:   make(chan bool, concurreny),
		task:    &wg,
		Feeds:   map[string]rssFeed{},
	}
}

func (r *Runner) SetSources(sources map[string]sourceInfo) {
	r.Sources = sources
}

func (r *Runner) DownloadFeeds(n int) {
	defer r.task.Wait()
	iterCount := 0
	for _, val := range r.Sources {
		iterCount++
		if iterCount > n {
			continue
		}
		r.task.Add(1)
		go r.getSourceFeed(val)
	}
}

func (r *Runner) RetryFailed() {
	defer r.task.Wait()
	log.Printf("Retrying failed: %d/%d", len(r.Failed), len(r.Sources))
	failed := r.Failed
	r.Failed = []string{}
	for _, key := range failed {
		source := r.Sources[key]
		r.setSchemeToHTTP(&source)
		r.task.Add(1)
		go r.getSourceFeed(source)
	}
}

func (r *Runner) addToFailed(source sourceInfo) {
	r.failedLock.Lock()
	r.Failed = append(r.Failed, source.Url)
	r.failedLock.Unlock()
}
func (r *Runner) addSourceFeed(source sourceInfo, feed rssFeed) {
	r.successLock.Lock()
	r.Feeds[source.Url] = feed
	r.successLock.Unlock()
}

func (r *Runner) setRequestHeaders(req *http.Request) {
	switch 1 {
	case 1:
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Host", "httpbin.org")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "cross-site")
		req.Header.Set("Sec-Fetch-User", "?1")
		// r.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv,90.0) Gecko/20100101 Firefox/90.0")
	}
}

func (r *Runner) getSourceFeed(source sourceInfo) {
	defer r.task.Done()
	r.limit <- true
	defer func() {
		<-r.limit
	}()
	resp, err := r.sendReq(source.Url)
	if err != nil {
		r.addToFailed(source)
		log.Println("failed: " + source.PublisherName)
		return
	}
	defer resp.Body.Close()
	respContent := rssFeed{}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%v", err)
		r.addToFailed(source)
		return
	} else {
		xml.Unmarshal(respBytes, &respContent)
		r.addSourceFeed(source, respContent)
		log.Println("success: " + source.PublisherName)
		return
	}

}

func (r *Runner) sendReq(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	r.setRequestHeaders(req)
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Runner) setSchemeToHTTP(source *sourceInfo) {
	httpUrl, _ := url.Parse(source.Url)
	httpUrl.Scheme = "http"
	source.Url = httpUrl.String()
}

func (r *Runner) Clean() map[string][]cleanedItem {
	cleaner := NewFeedParser(r.Feeds, r.Sources, time.Now().Add(10*time.Minute))
	articleMap := cleaner.ParseRss()
	return articleMap
}
