package models
import (
	"time"
	"github.com/google/uuid"
)
type Post struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Content      string    `json:"content" db:"content"`
	MediaURLs    []string  `json:"media_urls" db:"media_urls"`
	ScheduledAt  time.Time `json:"scheduled_at" db:"scheduled_at"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"`
	Status       string    `json:"status" db:"status"`
	PlatformPostID string  `json:"platform_post_id" db:"platform_post_id"`
	FailureReason string   `json:"failure_reason" db:"failure_reason"`
}
type PostRequest struct {
	Content          string    `json:"content" binding:"required"`
	ScheduledAt      time.Time `json:"scheduled_at" binding:"required"`
	TargetAccountIDs []uuid.UUID `json:"target_account_ids" binding:"required"`
}
