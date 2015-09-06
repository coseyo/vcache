package util

import (
	"fmt"
	"testing"
)

func TestSortedMaps(t *testing.T) {
	params := map[string]interface{}{
		"aa": 11,
		"cc": 22,
		"dd": "asdf",
		"bb": "dfe",
	}
	var m SortedMaps
	p := m.Sort(params)
	fmt.Println(p)
}
