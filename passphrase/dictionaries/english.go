package dictionaries

import (
	_ "embed" // for embedding dictionary files
	"sort"
	"strings"
	"sync"
)

//go:embed english.txt
var englishTxtRaw string

var (
	englishWords []string
	englishOnce  sync.Once
)

func init() {
	englishOnce.Do(func() {
		englishTxtRaw = strings.ReplaceAll(englishTxtRaw, "\r", "")
		englishWords = strings.Split(englishTxtRaw, "\n")
		sort.Strings(englishWords)
	})
}

// English returns all known English words.
func English() []string {
	rsp := make([]string, len(englishWords))
	for idx := range rsp {
		rsp[idx] = englishWords[idx]
	}
	return rsp
}
