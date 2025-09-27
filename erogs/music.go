package erogs

import (
	"encoding/json"
	"fmt"
  
	kurohelpererrors "kurohelper/errors"
	erogsmodels "kurohelper/models/erogs"
)

func GetMusicByFuzzy(search string) (*erogsmodels.FuzzySearchMusicResponse, error) {
	sql, err := buildFuzzySearchMusicSQL(search)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}
  
	var res erogsmodels.FuzzySearchMusicResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	if len(res.GameCategories) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	return &res, nil
}
