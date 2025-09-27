package vndb

import (
	"encoding/json"

	vndbmodels "kurohelper/models/vndb"
)

func GetStats() (*vndbmodels.Stats, error) {
	r, err := sendGetRequest("/stats")
	if err != nil {
		return nil, err
	}

	var res vndbmodels.Stats
	err = json.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
