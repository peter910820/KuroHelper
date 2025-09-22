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

func GetProducerByFuzzy(keyword string, companyType string) (*vndbmodels.ProducerSearchResponse, error) {
	reqProducer := vndbmodels.VndbCreate()

	filtersProducer := []interface{}{}
	if companyType != "" {
		filtersProducer = append(filtersProducer, "and")
		switch companyType {
		case "company":
			filtersProducer = append(filtersProducer, []string{"type", "=", "co"})
		case "individual":
			filtersProducer = append(filtersProducer, []string{"type", "=", "in"})
		case "amateur":
			filtersProducer = append(filtersProducer, []string{"type", "=", "ng"})
		}
		filtersProducer = append(filtersProducer, []string{"search", "=", keyword})
	} else {
		filtersProducer = []interface{}{"search", "=", keyword}
	}

	reqProducer.Filters = filtersProducer

	basicFields := "id, name, original, aliases, lang, type, description"
	extlinksFields := "extlinks.url, extlinks.label, extlinks.name, extlinks.id"

	allFields := []string{
		basicFields,
		extlinksFields,
	}

	reqProducer.Fields = strings.Join(allFields, ", ")

	jsonProducer, err := json.Marshal(reqProducer)
	if err != nil {
		return nil, err
	}

	respProducer, err := http.Post(os.Getenv("VNDB_ENDPOINT")+"/producer", "application/json", bytes.NewBuffer(jsonProducer))
	if err != nil {
		return nil, err
	}

	defer respProducer.Body.Close()

	r, err := io.ReadAll(respProducer.Body)
	if err != nil {
		return nil, err
	}

	if respProducer.StatusCode != 200 {
		return nil, fmt.Errorf("the server returned an error status code %d", respProducer.StatusCode)
	}

	var resProducer vndbmodels.BasicResponse[vndbmodels.ProducerSearchProducerResponse]
	err = json.Unmarshal(r, &resProducer)
	if err != nil {
		return nil, err
	}

	if len(resProducer.Results) == 0 {
		return nil, internalerrors.ErrVndbNoResult
	}

	// 等到查詢解析完後才能去查詢遊戲的資料
	reqVn := vndbmodels.VndbCreate()

	reqVn.Filters = []interface{}{
		"developer", "=", []interface{}{"id", "=", resProducer.Results[0].ID},
	}

	reqVn.Fields = "title, alttitle, length_minutes, length_votes, average, rating, votecount"

	jsonVn, err := json.Marshal(reqVn)
	if err != nil {
		return nil, err
	}
	respVn, err := http.Post(os.Getenv("VNDB_ENDPOINT")+"/vn", "application/json", bytes.NewBuffer(jsonVn))
	if err != nil {
		return nil, err
	}

	defer respVn.Body.Close()

	r, err = io.ReadAll(respVn.Body)
	if err != nil {
		return nil, err
	}

	if respVn.StatusCode != 200 {
		return nil, fmt.Errorf("the server returned an error status code %d", respVn.StatusCode)
	}

	var resVn vndbmodels.BasicResponse[vndbmodels.ProducerSearchVnResponse]
	err = json.Unmarshal(r, &resVn)
	if err != nil {
		return nil, err
	}

	if len(resVn.Results) == 0 {
		return nil, internalerrors.ErrVndbNoResult
	}

	return &vndbmodels.ProducerSearchResponse{
		Producer: resProducer,
		Vn:       resVn,
	}, nil
}
