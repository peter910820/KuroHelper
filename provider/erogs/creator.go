package erogs

import (
	"encoding/json"
	"fmt"
	"kurohelper/cache"
)

func GetCreatorByFuzzy(search string, idSearch bool) (*FuzzySearchCreatorResponse, error) {
	searchJP := ""
	if !idSearch {
		searchJP = cache.ZhTwToJp(search)
	}
	sql, err := buildFuzzySearchCreatorSQL(search, searchJP, idSearch)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res FuzzySearchCreatorResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func GetCreatorListByFuzzy(search string) (*[]FuzzySearchListResponse, error) {
	searchJP := cache.ZhTwToJp(search)
	sql, err := buildFuzzySearchCreatorListSQL(search, searchJP)
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
