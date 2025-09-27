package erogs

import (
	"encoding/json"

	kurohelpererrors "kurohelper/errors"
	erogsmodels "kurohelper/models/erogs"
)

func GetCreatorByFuzzy(search string) (*erogsmodels.FuzzySearchCreatorResponse, error) {
	sql, err := buildFuzzySearchCreatorSQL(search)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res erogsmodels.FuzzySearchCreatorResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	if len(res.Games) == 0 {
		return nil, kurohelpererrors.ErrVndbNoResult
	}

	return &res, nil
}
