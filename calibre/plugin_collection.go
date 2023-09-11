package calibre

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
)

func NewPluginCollection(json_file string) (Collection, error) {
	data, err := get_json_data(json_file)
	if err != nil {
		return nil, err
	}

	normalize_item := func(item string) (result string) {
		if strings.HasPrefix(item, "#") {
			result, _, _ = strings.Cut(item[1:], "^")
		} else {
			result = item
		}
		return
	}

	normalize_items := func(items []byte) (result []string) {
		jsonparser.ArrayEach(items, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if err == nil {
				result = append(result, normalize_item(string(value)))
			}
		})
		return
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
