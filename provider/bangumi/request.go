package bangumi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"os"

	kurohelpererrors "kurohelper/errors"

	"github.com/sirupsen/logrus"
)

func sendPostRequest(apiRoute string, jsonBytes []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, apiRoute, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	accessToken := os.Getenv("BANGUMI_ACCESS_TOKEN")
	userAgent := os.Getenv("BANGUMI_USER_AGENT")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logrus.Print(string(bodyBytes))
		return nil, fmt.Errorf("%w %d", kurohelpererrors.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func sendGetRequest(apiRoute string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, apiRoute, nil)
	if err != nil {
		return nil, err
	}
	accessToken := os.Getenv("BANGUMI_ACCESS_TOKEN")
	userAgent := os.Getenv("BANGUMI_USER_AGENT")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logrus.Print(string(bodyBytes))
		return nil, fmt.Errorf("%w %d", kurohelpererrors.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return r, nil
}
