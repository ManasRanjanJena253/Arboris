package api

import (
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
		ID int `json:"id"`
	} `json:"installation"`

	Repository struct {
		ID       int    `json:"id"`
		FullName string `json:"full_name"`
	} `json:"repository"`

	Sender struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
	} `json:"sender"`
}

type Job struct {
	InstallationID int    `json:"installation_id"`
	RepoFullName   string `json:"repo_full_name"`
	PRNumber       int    `json:"pr_number"`
	HeadSHA        string `json:"head_sha"`
}

var jobQueue = make(chan Job, 100) // TODO: MOCK JOB QUEUE

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
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

	switch eventType {
	case "ping":
		w.WriteHeader(http.StatusOK)
		return
	case "pull_request":
		var event PullRequestEvent
		err := json.Unmarshal(body, &event)

		if err != nil {
			http.Error(w, "Invalide payload", http.StatusBadRequest)
			return
		}

		switch event.Action {
		case "opened", "synchronize", "reopened":
			job := Job{
				InstallationID: event.Installation.ID,
				RepoFullName:   event.Repository.FullName,
				PRNumber:       event.PullRequest.Number,
				HeadSHA:        event.PullRequest.Head.SHA,
			}

			select {
			case jobQueue <- job: // TODO : Implement jobQueue
				w.WriteHeader(http.StatusOK)
			default:
				http.Error(w, "queue full", http.StatusServiceUnavailable)
			}
			return

		default:
			w.WriteHeader(http.StatusOK)
			return
		}
	default:
		w.WriteHeader(http.StatusOK)
		return
	}
}
