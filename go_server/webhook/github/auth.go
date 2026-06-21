package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type AuthUser struct {
	PemSecret []byte
	AppID     string
	Cache     *redis.Client
}

type InstallationTokenResponse struct {
	Token       string
	ExpiresAt   string
	Permissions []string
}

func (authDetails *AuthUser) GetInstallationToken(ctx context.Context, installId string) (string, error) {

	installToken, fetchErr := authDetails.Cache.Get(ctx, installId).Result()

	if fetchErr != nil && !errors.Is(fetchErr, redis.Nil) {
		slog.Error("Unable to get the pre-existing install token", "ERROR", fetchErr)
		return "", fetchErr
	}

	if installToken != "" {
		return installToken, nil
	}

	privateKey, pemErr := jwt.ParseRSAPrivateKeyFromPEM(authDetails.PemSecret)

	if pemErr != nil {
		slog.Error("Unable to parse the pem secret", "ERROR", pemErr)
		return "", pemErr
	}
	issuedAt := time.Now()
	claims := jwt.MapClaims{
		"iat": issuedAt.Unix(),
		"exp": issuedAt.Add(time.Second * 600).Unix(),
		"iss": authDetails.AppID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, signErr := token.SignedString(privateKey)

	if signErr != nil {
		slog.Error("Unable to sign the token", "ERROR", signErr)
		return "", signErr
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installId)
	req, reqErr := http.NewRequest("POST", url, nil)

	if reqErr != nil {
		slog.Error("Unable to create new request", "ERROR", reqErr)
		return "", reqErr
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-Github-Api-Version", "2022-11-28")

	client := &http.Client{Timeout: time.Second * 10}

	resp, respErr := client.Do(req)

	if respErr != nil {
		slog.Error("Unable to fetch installation id from github url", "ERROR", respErr)
		return "", respErr
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Unable to close the body", "ERROR", err)
			return
		}
	}(resp.Body)

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		slog.Error("Unable to read the response body", "ERROR", readErr)
		return "", readErr
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error from github api when fetching installation id", "Status Code", resp.StatusCode, "ERROR", resp.Status)
		return "", errors.New(fmt.Sprintf("API Side error. StatusCode : %d, Status : %s", resp.StatusCode, resp.Status))
	}

	var tokenResponse InstallationTokenResponse

	parseErr := json.Unmarshal(body, &tokenResponse)

	if parseErr != nil || tokenResponse.Token == "" {
		slog.Error("Unable to parse the response body", "ERROR", parseErr)
		return "", parseErr
	}

	exp, timeParseErr := time.Parse(time.RFC3339, tokenResponse.ExpiresAt)

	if timeParseErr != nil {
		slog.Error("Unable to parse expiry time provided by api", "Expiry Time", tokenResponse.ExpiresAt, "ERROR", timeParseErr)
		return "", timeParseErr
	}

	now := time.Now()
	setErr := authDetails.Cache.Set(ctx, installId, tokenResponse.Token, exp.Sub(now)).Err()

	if setErr != nil {
		slog.Info("Unable to set token in redis", "ERROR", setErr)
	}

	return tokenResponse.Token, nil
}
