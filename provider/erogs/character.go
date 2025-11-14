package erogs

import (
	"encoding/json"
	"fmt"
	"kurohelper/cache"
)

func GetCharacterByFuzzy(search string, idSearch bool) (*FuzzySearchCharacterResponse, error) {
	searchJP := ""
	if !idSearch {
		searchJP = cache.ZhTwToJp(search)
	}
	sql, err := buildFuzzySearchCharacterSQL(search, searchJP, idSearch)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res FuzzySearchCharacterResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func GetCharacterListByFuzzy(search string) (*[]FuzzySearchListResponse, error) {
	searchJP := cache.ZhTwToJp(search)
	sql, err := buildFuzzySearchCharacterListSQL(search, searchJP)
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
