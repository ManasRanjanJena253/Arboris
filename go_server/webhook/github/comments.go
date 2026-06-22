package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

type CommentDetailsResponse struct {
	ID        string `json:"id"`
	NodeID    string `json:"node_id"`
	Path      string `json:"path"`
	Side      string `json:"side"` // Can either be "RIGHT" or "LEFT"
	Line      int    `json:"line"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CommentReqBody struct {
	CommitID string `json:"commit_id"`
	Line     int    `json:"line"`
	Side     string `json:"side"` // Can either be "RIGHT" or "LEFT"
	Body     string `json:"body"`
	Path     string `json:"path"`
}

func (client *Client) PostComment(ctx context.Context, owner, repoName string, prNumber int, installID, comment, commitID, filePath, side string, line int) (*CommentDetailsResponse, error) {
	urlPath := fmt.Sprintf("/repos/%s/%s/pulls/%d/comments", owner, repoName, prNumber)

	var reqBody = CommentReqBody{
		CommitID: commitID,
		Line:     line,
		Side:     side,
		Body:     comment,
		Path:     filePath,
	}

	jsonReq, marshalErr := json.Marshal(reqBody)

	if marshalErr != nil {
		slog.Error("Unable to marshal comment request body", "ERROR", marshalErr)
		return nil, marshalErr
	}

	var result CommentDetailsResponse

	_, respErr := client.DoRequest(ctx, urlPath, "POST", installID, bytes.NewBuffer(jsonReq), &result)

	if respErr != nil {
		slog.Error("Unable to post comment", "ERROR", respErr)
		return nil, respErr
	}
	return &result, nil
}
