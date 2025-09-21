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
	"kurohelper/models"
)

func GetVnUseID(brandid string) (*models.VndbResponse[models.VndbGetVnUseIDResponse], error) {
	req := models.VndbCreate()

	req.Filters = []interface{}{
		"id", "=", brandid,
	}

	titleFields := "title, alttitle"
	imageFields := "image.url"
	developersFields := "developers.name, developers.original, developers.aliases"
	nameFields := "titles.lang, titles.title, titles.official, titles.main"
	staffFields := "staff.name, staff.role, staff.aliases.name"
	characterFields := "va.character.original, va.character.vns.role"
	lengthFields := "length_minutes, length_votes"
	scoreFields := "average, rating, votecount"
	relationsFields := "relations.titles.title"

	allFields := []string{
		titleFields,
		imageFields,
		developersFields,
		nameFields,
		staffFields,
		characterFields,
		lengthFields,
		scoreFields,
		relationsFields,
	}

	req.Fields = strings.Join(allFields, ", ")

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(os.Getenv("VNDB_ENDPOINT")+"/vn", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the server returned an error status code %d", resp.StatusCode)
	}

	var res models.VndbResponse[models.VndbGetVnUseIDResponse]
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, internalerrors.ErrVndbNoResult
	}

	return &res, nil
}

func GetVn() {

}
