package calibre

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/sandman7920/ccsync/config"
)

func NewMetaCollection(cfg config.MetaConfig) (Collection, error) {
	data, err := get_json_data(cfg.Path)
	if err != nil {
		return nil, err
	}

	var records Collection

	by_author := map[string][]string{}
	by_series := map[string][]string{}
	by_tags := map[string][]string{}

	parse_tags := func(data []byte) (result []string) {
		jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
			t, _ := jsonparser.ParseString(value)
			if cfg.ByTags == 1 {
				result = append(result, t)
			} else {
				v, found := cfg.Tags[strings.ToLower(t)]
				if found {
					result = append(result, v)
				}
			}
		})
		return
	}

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, _ error) {
		cdeKey, _ := jsonparser.GetString(value, "identifiers", "mobi-asin")
		if len(cdeKey) == 0 {
			cdeKey, _ = jsonparser.GetString(value, "uuid")
		}
		series, _ := jsonparser.GetString(value, "series")
		len_series := len(series)
		if len_series > 0 {
			by_author[series] = append(by_author[series], cdeKey)
		}

		if cfg.ByAuthor == 1 || (cfg.ByAuthor == 2 && len_series == 0) {
			author, _ := jsonparser.GetString(value, "authors", "[0]")
			if len(author) > 0 {
				by_author[author] = append(by_author[author], cdeKey)
			}
		}

		if cfg.ByTags > 0 {
			tags, _, _, err := jsonparser.Get(value, "tags")
			if err == nil {
				for _, t := range parse_tags(tags) {
					by_tags[t] = append(by_tags[t], cdeKey)
				}
			}
		}
	})

	for k, v := range by_author {
		records = append(records, Record{Title: k, CDEKeys: v})
	}

	for k, v := range by_series {
		records = append(records, Record{Title: k, CDEKeys: v})
	}

	for k, v := range by_tags {
		records = append(records, Record{Title: k, CDEKeys: v})
	}

	if err != nil {
		return nil, err
	}

	sort.SliceStable(records, func(lhs, rhs int) bool {
		return records[lhs].Title < records[rhs].Title
	})

	return records, nil
}
