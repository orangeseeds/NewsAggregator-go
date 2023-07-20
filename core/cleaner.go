package core

import (
	"errors"
	"time"
)

type feedParser struct {
	feeds      map[string]rssFeed
	sources    map[string]sourceInfo
	failedList []string
	getAfter   time.Time
}

func NewFeedParser(feeds map[string]rssFeed, sources map[string]sourceInfo, getAfter time.Time) *feedParser {
	return &feedParser{
		feeds:    feeds,
		sources:  sources,
		getAfter: getAfter,
	}
}

func (p *feedParser) beforeTime(t string) bool {
	parsedTime, err := p.parseTime(t)
	if err != nil || parsedTime.Before(p.getAfter) {
		return false
	}
	return true
}

func (p *feedParser) parseTime(t string) (*time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 02 Jan 2006 15:04:05 Z",
		"2006-01-02T15:04:05Z",
	}
	for i, format := range formats {
		parsedTime, err := time.Parse(format, t)
		if err == nil {
			return &parsedTime, nil
		} else if len(formats) == i+1 {
			return nil, err
		}
	}
	return nil, errors.New("time parse error: time format not provided")
}

func (p *feedParser) ParseRss() map[string][]cleanedItem {
	articleMap := map[string][]cleanedItem{}
	for key, val := range p.feeds {
		for _, article := range val.Channel.Items {
			if p.beforeTime(article.PubDate) {
				continue
			}
			if cleanedItem, ok := cleanItem(article, p.sources[key]); ok {
				articleMap[key] = append(articleMap[key], cleanedItem)
			}
		}
	}
	// log.Printf("Total error: %v out of %v", len(p.failedList), len(p.sources))
	return articleMap
}

func getCategory(item rssItem) []string {
	if item.Category == nil {
		return []string{}
	}
	return item.Category
}

func getMediaContents(item rssItem) []rssMedia {
	if item.MediaContents == nil {
		return []rssMedia{}
	}
	return item.MediaContents
}

func cleanItem(item rssItem, source sourceInfo) (cleanedItem, bool) {
	cleaned := cleanedItem{
		Title:         item.Title,
		Description:   item.Description,
		Link:          item.Link,
		PubDate:       item.PubDate,
		MediaContents: getMediaContents(item),
		Creator:       item.Creator,
		Category:      getCategory(item),
		Source: cleanedSource{
			Name:     source.PublisherName,
			ImageUrl: source.FaviconUrl,
			Url:      source.Url,
		},
	}
	return cleaned, true
}
