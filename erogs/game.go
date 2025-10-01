package erogs

import (
	"encoding/json"
	"fmt"

	kurohelpererrors "kurohelper/errors"
	erogsmodels "kurohelper/models/erogs"
)

func GetGameByFuzzy(search string, opt string) (*erogsmodels.FuzzySearchGameResponse, error) {
	if opt == "" {
		search = zhtwToJp(search)
	}
	sql, err := buildFuzzySearchGameSQL(search)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res erogsmodels.FuzzySearchGameResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	if len(res.CreatorShubetu) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	return &res, nil
}
