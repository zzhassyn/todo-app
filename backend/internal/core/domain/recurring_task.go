package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecurringTask struct {
	ID             int
	AuthorUserID   int
	Title          string
	Description    *string
	Priority       Priority
	FolderID       *uuid.UUID
	Tags           []string
	CronExpression string
	LastRunAt      *time.Time
	NextRunAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
