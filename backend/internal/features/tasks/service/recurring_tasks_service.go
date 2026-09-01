package tasks_service

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *TasksService) CreateRecurringTask(ctx context.Context, task domain.RecurringTask) (domain.RecurringTask, error) {
	// Validate cron
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := parser.Parse(task.CronExpression)
	if err != nil {
		return domain.RecurringTask{}, fmt.Errorf("invalid cron expression: %w", err)
	}

	task.NextRunAt = schedule.Next(time.Now().UTC())

	if err := s.checkFolderOwnership(ctx, task.FolderID, task.AuthorUserID); err != nil {
		return domain.RecurringTask{}, err
	}
	
	if len(task.Tags) > 0 {
		if err := validateTagNames(task.Tags); err != nil {
			return domain.RecurringTask{}, err
		}
	}

	created, err := s.tasksRepository.CreateRecurringTask(ctx, task)
	if err != nil {
		return domain.RecurringTask{}, fmt.Errorf("create recurring task: %w", err)
	}
	return created, nil
}

func (s *TasksService) GetRecurringTasks(ctx context.Context, userID int) ([]domain.RecurringTask, error) {
	tasks, err := s.tasksRepository.GetRecurringTasks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get recurring tasks: %w", err)
	}
	return tasks, nil
}

func (s *TasksService) DeleteRecurringTask(ctx context.Context, id int, userID int) error {
	if err := s.tasksRepository.DeleteRecurringTask(ctx, id, userID); err != nil {
		return fmt.Errorf("delete recurring task: %w", err)
	}
	return nil
}
