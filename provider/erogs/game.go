package erogs

import (
	"encoding/json"
	"fmt"
)

func GetGameByFuzzy(search string, idSearch bool) (*FuzzySearchGameResponse, error) {
	searchJP := ""
	if !idSearch {
		searchJP = zhtwToJp(search)
	}
	sql, err := buildFuzzySearchGameSQL(search, searchJP, idSearch)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res FuzzySearchGameResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

func GetGameListByFuzzy(search string) (*[]FuzzySearchListResponse, error) {
	searchJP := zhtwToJp(search)
	sql, err := buildFuzzySearchGameListSQL(search, searchJP)
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
