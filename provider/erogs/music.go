package erogs

import (
	"encoding/json"
	"fmt"
)

func GetMusicByFuzzy(search string, opt string) (*FuzzySearchMusicResponse, error) {
	if opt == "" {
		search = zhtwToJp(search)
	}
	sql, err := buildFuzzySearchMusicSQL(search)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res FuzzySearchMusicResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}
