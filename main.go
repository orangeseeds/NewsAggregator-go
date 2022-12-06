package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/orangeseeds/aggregator-go/core"
)

func fillNewFileWith(data interface{}) {
	fileName := fmt.Sprintf("./results/rss-%v.json", time.Now().Unix())
	log.Println("Data written to " + fileName)
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	feedJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	f.Write(feedJson)
}

func main() {
	sourceMap := core.LoadSources("feed.json")
	// core.GetAfter = time.Now().Add(- * time.Hour)
	runner := core.NewRunner(64, sourceMap)
	runner.DownloadFeeds(len(sourceMap))
	runner.RetryFailed()
	log.Printf("Failed: %d", len(runner.Failed))
	articles := runner.Clean()

	fillNewFileWith(articles)

}
