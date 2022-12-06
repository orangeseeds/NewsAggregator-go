package core

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
)

var (
// SourceMap  = map[string]sourceInfo{}
// task       sync.WaitGroup
// failMutex  sync.Mutex
// failedList []sourceInfo
// client     http.Client
// FeedMap    = map[string]rssFeed{}
// feedMutex  sync.Mutex
// GetAfter time.Time = time.Now().Add(-24 * time.Hour)
// limit              = make(chan bool, 5)
)

type sourceInfo struct {
	Category           string   `json:"category"`
	PublisherName      string   `json:"publisher_name"`
	ContentType        string   `json:"content_type"`
	PublisherDomain    string   `json:"publisher_domain"`
	MaxEntries         int      `json:"max_entries"`
	Url                string   `json:"url"`
	CoverUrl           string   `json:"cover_url"`
	FaviconUrl         string   `json:"favicon_url"`
	BackgroundColor    string   `json:"background_color"`
	DestinationDomains string   `json:"destination_domains"`
	Channels           []string `json:"channels"`
}

type rssFeed struct {
	XMLNAME xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssMedia struct {
	URl    string `xml:"url,attr"`
	Medium string `xml:"medium,attr"`
}

type rssItem struct {
	Title         string     `xml:"title"`
	Description   string     `xml:"description"`
	Link          string     `xml:"link"`
	PubDate       string     `xml:"pubDate"`
	MediaContents []rssMedia `xml:"content"`
	Creator       string     `xml:"creator"`
	Category      []string   `xml:"category"`
}

type cleanedSource struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	ImageUrl string `json:"image_url"`
}
type cleanedItem struct {
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	Link          string        `json:"link"`
	PubDate       string        `json:"pub_date"`
	MediaContents []rssMedia    `json:"media_contents"`
	Creator       string        `json:"creator"`
	Category      []string      `json:"category"`
	Source        cleanedSource `json:"source"`
}

type rssChannel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Image       struct {
		URL string `xml:"url"`
	} `xml:"image"`
	Language string    `xml:"language"`
	Items    []rssItem `xml:"item"`
}

func LoadSources(file string) map[string]sourceInfo {
	sourceMap := map[string]sourceInfo{}
	jsonFeed, _ := ioutil.ReadFile(file)
	err := json.Unmarshal(jsonFeed, &sourceMap)
	if err != nil {
		panic(err)
	}
	return sourceMap
}
