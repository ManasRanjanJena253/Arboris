package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type contextKey string

const InstallationIDKey contextKey = "installationID"

func ExtractInstallID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read the payload", http.StatusInternalServerError)
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		var payload struct {
			Installation struct {
				ID int64 `json:"id"`
			} `json:"installation"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), InstallationIDKey, payload.Installation.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
