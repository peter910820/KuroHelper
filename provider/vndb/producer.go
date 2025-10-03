package vndb

import (
	"encoding/json"
	"strings"

	kurohelpererrors "kurohelper/errors"
)

func GetProducerByFuzzy(keyword string, companyType string) (*ProducerSearchResponse, error) {
	reqProducer := VndbCreate()

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

	r, err := sendPostRequest("/producer", jsonProducer)
	if err != nil {
		return nil, err
	}

	var resProducer BasicResponse[ProducerSearchProducerResponse]
	err = json.Unmarshal(r, &resProducer)
	if err != nil {
		return nil, err
	}

	if len(resProducer.Results) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	// 等到查詢解析完後才能去查詢遊戲的資料
	reqVn := VndbCreate()

	reqVn.Filters = []interface{}{
		"developer", "=", []interface{}{"id", "=", resProducer.Results[0].ID},
	}

	reqVn.Fields = "title, alttitle, length_minutes, length_votes, average, rating, votecount"

	jsonVn, err := json.Marshal(reqVn)
	if err != nil {
		return nil, err
	}

	r, err = sendPostRequest("/vn", jsonVn)
	if err != nil {
		return nil, err
	}

	var resVn BasicResponse[ProducerSearchVnResponse]
	err = json.Unmarshal(r, &resVn)
	if err != nil {
		return nil, err
	}

	if len(resVn.Results) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	return &ProducerSearchResponse{
		Producer: resProducer,
		Vn:       resVn,
	}, nil
}
