package api

import (
	"Arboris/go_server/webhook/middleware"
	"Arboris/go_server/webhook/queue"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type PullRequestEvent struct {
	Action string `json:"action"`

	PullRequest struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		State  string `json:"state"`
		URL    string `json:"url"`
		Head   struct {
			SHA string `json:"sha"`
		} `json:"head"`
	} `json:"pull_request"`

	Installation struct {
		ID string `json:"id"`
	} `json:"installation"`

	Repository struct {
		ID       int    `json:"id"`
		FullName string `json:"full_name"`
		Owner    struct {
			Login string `json:"login"`
			ID    int    `json:"id"`
			Type  string `json:"type"`
		}
	} `json:"repository"`

	Sender struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
	} `json:"sender"`
}
type WebhookHandlerStruct struct {
	queue *queue.JobQueue
}

func (hookHandler *WebhookHandlerStruct) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	eventType := r.Header.Get("X-Github-Event")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Unable to read Payload", http.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Unable to close the request body", "ERROR", err)
		}
	}(r.Body)

	installId, ok := r.Context().Value(middleware.InstallationIDKey).(int64)

	if !ok {
		slog.Error("The installId doesn't exist")
		http.Error(w, "InstallID Not found", http.StatusBadRequest)
		return
	}

	switch eventType {
	case "ping":
		w.WriteHeader(http.StatusOK)
		return
	case "pull_request":
		var event PullRequestEvent
		err := json.Unmarshal(body, &event)

		if err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		switch event.Action {
		case "opened", "synchronize", "reopened":
			job := queue.Job{
				Owner:          event.Repository.Owner.Login,
				InstallationID: installId,
				RepoName:       event.Repository.FullName,
				PRNumber:       event.PullRequest.Number,
				CommitID:       event.PullRequest.Head.SHA,
			}

			enqueueErr := hookHandler.queue.JobEnqueue(&job)
			if !enqueueErr {
				http.Error(w, "Queue Full", http.StatusServiceUnavailable)
				return
			}

		default:
			w.WriteHeader(http.StatusOK)
			return
		}
	default:
		w.WriteHeader(http.StatusOK)
		return
	}
}
