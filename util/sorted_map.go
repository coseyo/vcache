package util

import "sort"

type sortedMap struct {
	Key   string
	Value interface{}
}

type SortedMaps []sortedMap

func (this SortedMaps) Len() int {
	return len(this)
}

func (this SortedMaps) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this SortedMaps) Less(i, j int) bool {
	return this[i].Key < this[j].Key
}

func (this SortedMaps) Sort(params map[string]interface{}) map[string]interface{} {
	for k, v := range params {
		this = append(this, sortedMap{Key: k, Value: v})
	}
	sort.Sort(this)
	m := map[string]interface{}{}
	for _, v := range this {
		m[v.Key] = v.Value
	}
	return m
}
