package kindle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Commander struct{}

var client = &http.Client{
	Timeout: 2 * time.Minute,
}

var token = func() string {
	fp, err := os.Open("/tmp/session_token")
	if err != nil {
		return ""
	}
	defer fp.Close()
	data, err := io.ReadAll(fp)
	if err != nil {
		return ""
	}

	return strings.Trim(string(data), "\r\n\t ")
}()

func send_cmd(json_data []byte) error {
	req, err := http.NewRequest("POST", "http://127.0.0.1:9101/change", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Conent-Type", "application/json")
	if len(token) > 0 {
		// workaround golang CanonicalHeaderKey
		req.Header["AuthToken"] = []string{token}
	}
	req.Close = true

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("send_cmd: invalid status_code: %d, expected: 200", res.StatusCode)
	}
	// fmt.Print(io.ReadAll(res.Body))

	return err
}

func (u *Commander) Purge(collection Collection) error {
	cmd := make([]map[string]interface{}, 0)
	for _, c := range collection {
		if len(c.UUID) > 0 {
			cmd = append(cmd, map[string]interface{}{
				"delete": map[string]interface{}{
					"uuid": c.UUID,
				},
			})
		}
	}

	update := map[string]interface{}{
		"commands": cmd,
		"type":     "ChangeRequest",
		"id":       1,
	}

	json_data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	return send_cmd(json_data)
}

func (u *Commander) Update(collection Collection, isCcAware bool, unique map[string]int) error {
	ts := time.Now().Unix()
	cmd := make([]map[string]interface{}, 0)
	for _, c := range collection {
		if len(c.UUID) > 0 {
			cmd = append(cmd, map[string]interface{}{
				"delete": map[string]interface{}{
					"uuid": c.UUID,
				},
			})
		}
		uuid := uuid.New()
		cmd = append(cmd, map[string]interface{}{
			"insert": map[string]interface{}{
				"type":       "Collection",
				"uuid":       uuid,
				"lastAccess": ts,
				"titles": []map[string]interface{}{{
					"display": c.Title,
				}},
				"isVisibleInHome":       1,
				"isArchived":            1,
				"mimeType":              "application/x-kindle-collection",
				"collections":           nil,
				"collectionCount":       nil,
				"collectionDataSetName": uuid,
			},
			"update": map[string]interface{}{
				"type":    "Collection",
				"uuid":    uuid,
				"members": c.Members(),
			},
		})
	}

	if isCcAware {
		for uuid, count := range unique {
			cmd = append(cmd, map[string]interface{}{
				"update": map[string]interface{}{
					"type":            "Entry:Item",
					"uuid":            uuid,
					"collectionCount": count,
				},
			})
		}
	}

	update := map[string]interface{}{
		"commands": cmd,
		"type":     "ChangeRequest",
		"id":       1,
	}

	json_data, err := json.Marshal(update)
	if err != nil {
		return err
	}
	return send_cmd(json_data)
}
