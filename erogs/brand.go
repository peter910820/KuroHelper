package erogs

import (
	"encoding/json"
	"fmt"

	kurohelpererrors "kurohelper/errors"
	erogsmodels "kurohelper/models/erogs"
)

func GetBrand(search string) (*erogsmodels.SearchBrandResponse, error) {
	sql, err := buildSearchBrandSQL(search)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res erogsmodels.SearchBrandResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	if len(res.GameList) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	return &res, nil
}
