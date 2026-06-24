package queue

import (
	"Arboris/go_server/webhook/github"
	"context"
	"log/slog"
	"time"
)

func StartWorkers(q *jobQueue, githubClient *github.Client, count int) {
	for i := 0; i < count; i++ {
		go worker(q, githubClient)
	}
}

func worker(q *jobQueue, githubClient *github.Client) {
	for job := range q.queue {
		proces(job, githubClient)
	}
}

func proces(job *Job, githubClient *github.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	changedFiles, changeErr := githubClient.GetChanges(ctx, job.Owner, job.RepoName, job.PRNumber, job.InstallationID)

	if changeErr != nil {
		slog.Error("Failed to fetch changes", "RepoName", job.RepoName, "PRNumber", job.PRNumber, "ERROR", changeErr)
		return
	}

	if len(changedFiles) == 0 {
		slog.Info("No changed files found", "pr", job.PRNumber)
		return
	}

	_, postErr := githubClient.PostComment(
		ctx,
		job.Owner,
		job.RepoName,
		job.PRNumber,
		job.InstallationID,
		"Arboris is analyzing this PR...",
		job.CommitID,
		changedFiles[0].Filename, // TODO: Use all the file lines and changes. Currently just a placeholder.
		"RIGHT",
		1,
	)
	if postErr != nil {
		slog.Error("Failed to post comment", "ERROR", postErr, "pr", job.PRNumber)
	}
}
