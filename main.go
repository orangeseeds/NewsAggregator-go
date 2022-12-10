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

	ticker := time.NewTicker(10 * time.Minute)

	runAggregator := func() {
		sourceMap := core.LoadSources("feed.json")
		runner := core.NewRunner(64, sourceMap)
		runner.DownloadFeeds(len(sourceMap))
		runner.RetryFailed()
		log.Printf("Failed: %d", len(runner.Failed))
		articles := runner.Clean()
		fillNewFileWith(articles)
	}

	log.Printf("Iteration: %d", 0)
	runAggregator()
loop:
	for i := 1; ; i++ {
		<-ticker.C
		log.Printf("Iteration: %d", i)
		runAggregator()
		if i == 10 {
			break loop
		}
	}

}

// t := time.NewTicker(1 * time.Second)
// // runChan := make(chan bool, 1)

// signalChan := make(chan os.Signal, 1)
// signal.Notify(signalChan, os.Interrupt)

// go func() {
// 	for i := 6; i > 0; i-- {
// 		<-t.C
// 		log.Printf("Exiting after %d secs..", i)
// 	}
// 	// <-runChan
// }()

// smoothExit := func() {
// 	sig := <-signalChan
// 	t.Stop()
// 	// close(runChan)
// 	log.Println("Closing the program in 2 secs...: ", sig)
// 	time.Sleep(2 * time.Second)
// }
// smoothExit()

// go smoothExit()
