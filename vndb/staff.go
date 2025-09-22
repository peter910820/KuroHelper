package vndb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	internalerrors "kurohelper/errors"
	vndbmodels "kurohelper/models/vndb"
)

func GetStaffByFuzzy(keyword string, roleType string) (*vndbmodels.BasicResponse[vndbmodels.StaffSearchResponse], error) {
	req := vndbmodels.VndbCreate()

	filters := []interface{}{}
	if roleType != "" {
		filters = append(filters, "and")
		// 傳進來的直接就是API篩選項規格的字串
		filters = append(filters, []string{"type", "=", roleType})
		filters = append(filters, []string{"search", "=", keyword})
	} else {
		filters = []interface{}{"search", "=", keyword}
	}

	req.Filters = filters

	basicFields := "id, aid, ismain, name, original, lang, gender, description"
	extlinksFields := "extlinks{url, label, name, id}"
	aliasesFields := "aliases{aid, name, latin, ismain}"

	allFields := []string{
		basicFields,
		extlinksFields,
		aliasesFields,
	}

	req.Fields = strings.Join(allFields, ", ")

	jsonProducer, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(os.Getenv("VNDB_ENDPOINT")+"/staff", "application/json", bytes.NewBuffer(jsonProducer))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the server returned an error status code %d", resp.StatusCode)
	}

	var res vndbmodels.BasicResponse[vndbmodels.StaffSearchResponse]
	err = json.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, internalerrors.ErrVndbNoResult
	}

	return &res, nil
}
