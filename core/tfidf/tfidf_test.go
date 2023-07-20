package tfidf

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestTfidf(t *testing.T) {
	testText := `Placed at the top of the device. When the main chamber fills up with the drug from PDMS micro-channel. this pushes the the membrane generating an electric signal. this electric signal triggres the actuator which then vibrates the diaphragm. this then creates a positive/negative volume in the chamber. The inlet valve and outlet valve controls the direction of flow.`

	stopWords, err := LoadStopWords("./stop_words.json")
	if err != nil {
		t.Fatal(err)
	}
	var validTokens [][]string
	tokens := LoadText(testText)
	for _, s := range tokens {
		var validSent []string
		for _, w := range s {
			w = strings.ToLower(w)
			if slices.Contains(stopWords, w) {
				continue
			}
			validSent = append(validSent, w)
		}
		validTokens = append(validTokens, validSent)
	}

	start := currentTimeMillis()
	tfidfs := FeatureSelect(validTokens)
	for _, item := range tfidfs {
		t.Log(item)
	}

	cost := currentTimeMillis() - start
	fmt.Printf("time consuming %d ms ", cost)
}