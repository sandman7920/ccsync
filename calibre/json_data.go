package calibre

import (
	"io"
	"os"
)

func get_json_data(json_file string) ([]byte, error) {
	fp, err := os.Open(json_file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return io.ReadAll(fp)
}
