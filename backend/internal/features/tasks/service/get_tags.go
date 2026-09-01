package tasks_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *TasksService) GetTags(ctx context.Context) ([]domain.Tag, error) {
	tags, err := s.tasksRepository.GetTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tags from repository: %w", err)
	}
	return tags, nil
}
