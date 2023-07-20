package tfidf

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"sort"
	"strings"
	"time"
)

type wordTfIdf struct {
	nworld string
	value  float64
}

type wordTfIdfs []wordTfIdf
type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func (us wordTfIdfs) Len() int {
	return len(us)
}
func (us wordTfIdfs) Less(i, j int) bool {
	return us[i].value > us[j].value
}
func (us wordTfIdfs) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
func FeatureSelect(list_words [][]string) wordTfIdfs {
	docFrequency := make(map[string]float64, 0)
	sumWorlds := 0
	for _, wordList := range list_words {
		for _, v := range wordList {
			docFrequency[v] += 1
			sumWorlds++
		}
	}
	wordTf := make(map[string]float64)
	for k := range docFrequency {
		wordTf[k] = docFrequency[k] / float64(sumWorlds)
	}
	docNum := float64(len(list_words))
	wordIdf := make(map[string]float64)
	wordDoc := make(map[string]float64, 0)
	for k := range docFrequency {
		for _, v := range list_words {
			for _, vs := range v {
				if k == vs {
					wordDoc[k] += 1
					break
				}
			}
		}
	}
	for k := range docFrequency {
		wordIdf[k] = math.Log(docNum / (wordDoc[k] + 1))
	}
	var wordifS wordTfIdfs
	for k := range docFrequency {
		var wti wordTfIdf
		wti.nworld = k
		wti.value = wordTf[k] * wordIdf[k]
		wordifS = append(wordifS, wti)
	}
	sort.Sort(wordifS)
	// fmt.Println(wordifS)
	return wordifS
}

func Load() [][]string {
	slice := [][]string{
		{"my", "dog", "has", "flea", "problems", "help", "please"},
		{"maybe", "not", "take", "him", "to", "dog", "park", "stupid"},
		{"my", "dalmation", "is", "so", "cute", "I", "love", "him"},
		{"stop", "posting", "stupid", "worthless", "garbage"},
		{"mr", "licks", "ate", "my", "steak", "how", "to", "stop", "him"},
		{"quit", "buying", "worthless", "dog", "food", "stupid"},
		{"a", "language", "can", "be", "significantly", "reduced", "in", "a", "good", "programming", "environment"},
		{"Programming", "environments", "are", "discussed", "in", "Section"},
		{"Fourth", "the", "cost", "of", "executing", "programs", "written", "in", "a", "language", "is", "greatly", "influenced", "by", "that", "design"},
		{"A", "language", "that", "requires", "many", "time", "type", "checks", "will", "prohibit", "fast", "code", "execution", "regardless", "of", "the", "quality", "of", "the", "compiler"},
		{"Although", "execution", "efficiency", "was", "the", "foremost", "concern", "in", "the", "design", "of", "early", "languages", "it", "is", "now", "considered", "to", "be", "less", "important"},
	}
	return slice
}

func LoadText(text string) [][]string {
	result := [][]string{}
	sentences := strings.Split(text, ".")
	for _, sentence := range sentences {
		words := strings.Split(sentence, " ")
		result = append(result, words)
	}
	return result
}

func LoadStopWords(path string) ([]string, error) {
	var stopWords = map[string][]string{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &stopWords)
	if err != nil {
		return nil, err
	}

	wordList, ok := stopWords["stop_words"]
	if !ok {
		return nil, err
	}

	return wordList, nil
}
