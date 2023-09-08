package calibre

import (
	"io"
	"os"
	"sort"
	"strings"

	"github.com/buger/jsonparser"
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

func get_json_data(json_file string) ([]byte, error) {
	fp, err := os.Open(json_file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return io.ReadAll(fp)
}

func normalize_item(item string) (result string) {
	if strings.HasPrefix(item, "#") {
		result, _, _ = strings.Cut(item[1:], "^")
	} else {
		result = item
	}
	return
}

func normalize_items(items []byte) (result []string) {
	jsonparser.ArrayEach(items, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err == nil {
			result = append(result, normalize_item(string(value)))
		}
	})
	return
}

func NewCollection(json_file string) (Collection, error) {
	data, err := get_json_data(json_file)
	if err != nil {
		return nil, err
	}

	var records Collection

	err = jsonparser.ObjectEach(data, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
		if idx := strings.LastIndex(string(key), "@"); idx > -1 {
			key = key[:idx]
		}
		items, _, _, err := jsonparser.Get(value, "items")
		if err != nil {
			return err
		}

		records = append(records, Record{
			Title:   string(key),
			CDEKeys: normalize_items(items),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.SliceStable(records, func(lhs, rhs int) bool {
		return records[lhs].Title < records[rhs].Title
	})

	return records, nil
}
