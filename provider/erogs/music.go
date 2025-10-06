package erogs

import (
	"encoding/json"
	"fmt"
)

func GetMusicByFuzzy(search string) (*FuzzySearchMusicResponse, error) {
	searchJP := zhtwToJp(search)
	sql, err := buildFuzzySearchMusicSQL(search, searchJP)
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

func GetMusicListByFuzzy(search string) (*[]FuzzySearchListResponse, error) {
	searchJP := zhtwToJp(search)
	sql, err := buildFuzzySearchMusicListSQL(search, searchJP)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res []FuzzySearchListResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}
