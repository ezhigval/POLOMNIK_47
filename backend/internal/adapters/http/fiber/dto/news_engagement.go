package dto

import (
	"time"

	"palomnik/internal/domain"
)

type NewsLikeStateResponse struct {
	LikeCount  int  `json:"like_count"`
	LikedByYou bool `json:"liked_by_you"`
}

type NewsCommentResponse struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type AddNewsCommentRequest struct {
	Body string `json:"body"`
}

func ToNewsLikeStateResponse(count int, liked bool) NewsLikeStateResponse {
	return NewsLikeStateResponse{LikeCount: count, LikedByYou: liked}
}

func ToNewsCommentResponse(comment domain.NewsComment) NewsCommentResponse {
	return NewsCommentResponse{
		ID:        comment.ID.String(),
		Author:    comment.Author,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt.UTC().Format(time.RFC3339),
	}
}
