package ymgal

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	kurohelpererrors "kurohelper/errors"

	"github.com/sirupsen/logrus"
)

// 做一次重試(取新Token)的版本
func sendWithRetry(apiRoute string) ([]byte, error) {
	r, err := sendGetRequest(apiRoute)
	if err != nil {
		if errors.Is(err, kurohelpererrors.ErrYmgalInvalidAccessToken) {
			logrus.Warnf("%s, refreshing and retrying...", err)
			err = GetToken()
			if err != nil {
				return nil, err
			}

			r, err = sendGetRequest(apiRoute)
			if err != nil {
				return nil, err
			}

			return r, nil
		}
		return nil, err
	}
	return r, nil
}

func sendGetRequest(apiRoute string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, apiRoute, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("version", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, kurohelpererrors.ErrYmgalInvalidAccessToken
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w %d", kurohelpererrors.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return r, nil
}
