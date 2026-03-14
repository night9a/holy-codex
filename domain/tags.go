package domain

import (
	"sort"
	"strings"
)

type TagManager struct {
	tags map[string]int
}

func NewTagManager() *TagManager {
	return &TagManager{tags: make(map[string]int)}
}

func (tm *TagManager) Add(tags []string) {
	for _,t := range tags {
		key := normalize(t)
		if key != "" {
			tm.tags[key]++
		}
	}
}
func (tm *TagManager) Remove(tags []string) {
	for _,t := range tags {
		key := normalize(t)
		if count,ok := tm.tags[key]; ok {
			if count <= 1 {
				delete(tm.tags,key)
			} else {
				tm.tags[key]--
			}
		}
	}
}


func (tm *TagManager) All() []string {
	list := make([]string,0,len(tm.tags))
	for t := range tm.tags {
		list = append(list, t)
	}
	sort.Strings(list)
	return list
}

func (tm *TagManager) Popular(n int) []string {
	type kv struct {
		tag string
		count int
	}
	pairs := make([]kv,0,len(tm.tags))
	for t,c := range tm.tags {
		pairs = append(pairs, kv{t,c})
	}
	sort.Slice(pairs, func(i,j int) bool {
		return pairs[i].count > pairs[j].count
	})
	result := make([]string,0,n)
	for i := 0;i < n && i < len(pairs); i++ {
		result = append(result, pairs[i].tag)
	}
	return result
}

func normalize(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}