package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(s string) []string {
	freqMap := make(map[string]int)

	for _, v := range strings.Fields(s) {
		freqMap[v]++
	}

	sortedSlice := make([]string, 0, len(freqMap))
	for k := range freqMap {
		sortedSlice = append(sortedSlice, k)
	}

	sort.Slice(sortedSlice, func(i, j int) bool {
		if freqMap[sortedSlice[i]] == freqMap[sortedSlice[j]] {
			return sortedSlice[i] < sortedSlice[j]
		}
		return freqMap[sortedSlice[i]] > freqMap[sortedSlice[j]]
	})

	if len(sortedSlice) < 10 {
		return sortedSlice
	}
	return sortedSlice[:10]
}
