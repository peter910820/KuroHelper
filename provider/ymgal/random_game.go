package ymgal

import (
	"encoding/json"
	"os"
)

func GetRandomGame() ([]randomGameResp, error) {
	r, err := sendWithRetry(os.Getenv("YMGAL_ENDPOINT") + "/open/archive/random-game?num=1")
	if err != nil {
		return nil, err
	}

	var game basicResp[[]randomGameResp]

	err = json.Unmarshal(r, &game)
	if err != nil {
		return nil, err
	}

	return game.Data, nil
}
