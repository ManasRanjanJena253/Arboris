package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Client struct {
	Client  *http.Client
	BaseURL string
	Auth    *AuthUser
}

func (client *Client) DoRequest(ctx context.Context, path string, method string, installID string, body io.Reader, responseStruct interface{}) ([]byte, error) {
	installToken, tokenErr := client.Auth.GetInstallationToken(ctx, installID)

	if tokenErr != nil {
		return nil, tokenErr
	}

	urlPath := client.BaseURL + path
	req, reqErr := http.NewRequest(method, urlPath, body)

	if reqErr != nil {
		slog.Error("Unable to create new http request", "ERROR", reqErr)
		return nil, reqErr
	}

	req.Header.Add("Authorization", "Bearer "+installToken)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-Github-Api-Version", "2022-11-28")

	resp, respErr := client.Client.Do(req)

	if respErr != nil {
		slog.Error(fmt.Sprintf("Unable to get response from the url : %s", urlPath))
		return nil, respErr
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Info("Unable to close the response body", "ERROR", err)
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		slog.Error("Error from github api when fetching installation id", "Status Code", resp.StatusCode, "ERROR", resp.Status)
		return nil, errors.New(fmt.Sprintf("API Side error. StatusCode : %d, Status : %s", resp.StatusCode, resp.Status))
	}

	respBody, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		slog.Error("Unable to read the response body", "ERROR", readErr)
		return nil, readErr
	}

	if len(respBody) == 0 {
		return nil, nil
	}

	parseErr := json.Unmarshal(respBody, responseStruct)

	if parseErr != nil {
		slog.Error("Unable to parse the response with given response struct", "ERROR", parseErr)
		return nil, parseErr
	}

	return respBody, nil
}
