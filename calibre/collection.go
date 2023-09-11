package calibre

import (
	"sort"
	"strings"
)

type Record struct {
	CDEKeys []string
	Title   string
}

type Collection []Record

func (c Collection) IndexOf(title string) int {
	idx, found := sort.Find(len(c), func(i int) int {
		return strings.Compare(title, c[i].Title)
	})

	if found {
		return idx
	}

	return -1
}

func (c Collection) Record(title string) *Record {
	idx := c.IndexOf(title)
	if idx == -1 {
		return nil
	}

	return &c[idx]
}
