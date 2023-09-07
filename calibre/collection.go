package calibre

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"strings"
)

type Entity struct {
	CDEKey  string
	CdeType string
}

type Record struct {
	Entities []Entity
	Title    string
}

func (r *Record) CDEKeys() []string {
	result := make([]string, 0, len(r.Entities))
	for _, e := range r.Entities {
		result = append(result, e.CDEKey)
	}
	return result
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

// Returns *calibre.Record
//
// if record is not found returns nil
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

func normalize_item(item string) (result Entity) {
	if strings.HasPrefix(item, "#") {
		result.CDEKey, result.CdeType, _ = strings.Cut(item[1:], "^")
	} else {
		result.CDEKey = item
		result.CdeType = "EBOK"
	}
	return
}

func normalize_items(values map[string]interface{}) (result []Entity) {
	items, found := values["items"]
	if !found {
		return
	}
	for _, item := range items.([]interface{}) {
		result = append(result, normalize_item(item.(string)))
	}
	return
}

func NewCollection(json_file string) (Collection, error) {
	data, err := get_json_data(json_file)
	if err != nil {
		return nil, err
	}

	var records Collection
	var cc map[string]map[string]interface{}
	err = json.Unmarshal(data, &cc)
	if err != nil {
		return nil, err
	}
	for k, v := range cc {
		if idx := strings.LastIndex(k, "@"); idx > -1 {
			k = k[:idx]
		}

		records = append(records, Record{
			Title:    k,
			Entities: normalize_items(v),
		})
	}

	sort.SliceStable(records, func(lhs, rhs int) bool {
		return records[lhs].Title < records[rhs].Title
	})

	return records, nil
}
