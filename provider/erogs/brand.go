package erogs

import (
	"encoding/json"
	"fmt"
)

func GetBrandBYFuzzy(search string, opt string) (*SearchBrandResponse, error) {
	sql, err := buildSearchBrandSQL(search)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res SearchBrandResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}
