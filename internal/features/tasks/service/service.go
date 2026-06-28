package tasks_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
	usersChecker    UsersChecker
	foldersChecker  FoldersChecker
}

// TasksFilter is a service-level (domain-oriented) filter for listing tasks.
// It is intentionally decoupled from any repository implementation detail.
//
// Archived defaults to "exclude archived tasks" when nil — see
// get_tasks.go's GetTasks for the exact semantics (nil => only
// non-archived, false => only non-archived, true => only archived).
//
// Folder filtering has three states, not two: leave both NoFolder and
// FolderID unset to not filter by folder at all; set NoFolder=true to see
// only unfiled tasks (the default "buffer"); set FolderID to see only
// tasks in that folder. NoFolder and FolderID are mutually exclusive —
// the transport layer is responsible for only ever setting one.
type TasksFilter struct {
	AuthorUserID *int
	Completed    *bool
	Archived     *bool
	Priority     *domain.Priority
	Tag          *string
	FolderID     *uuid.UUID
	NoFolder     bool
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task, tagNames []string) (domain.Task, error)
	GetTasks(
		ctx context.Context,
		filter TasksFilter,
		limit *int,
		offset *int,
	) ([]domain.Task, error)
	GetTask(ctx context.Context, id int) (domain.Task, error)
	ArchiveTask(ctx context.Context, id int, task domain.Task) (domain.Task, error)
	PatchTask(ctx context.Context, id int, task domain.Task, tagNames []string) (domain.Task, error)
	GetTags(ctx context.Context) ([]domain.Tag, error)
	// PermanentlyDeleteTask hard-deletes a task row. Only called by the
	// service after confirming the task is already archived — see
	// PermanentlyDeleteTask in permanently_delete_task.go.
	PermanentlyDeleteTask(ctx context.Context, id int) error
}

// UsersChecker is the subset of the users feature's service that the tasks
// feature depends on. The tasks feature does not depend on the users
// repository directly; it depends on this narrow interface, which is
// satisfied by users_service.UsersService.
type UsersChecker interface {
	GetUser(ctx context.Context, id int) (domain.User, error)
}

// FoldersChecker is the subset of the folders feature's service that the
// tasks feature depends on, used to verify a folder_id supplied on
// create/patch both exists and belongs to the requesting user before the
// task is written. Satisfied by folders_service.FoldersService. The tasks
// feature does not import the folders repository directly.
type FoldersChecker interface {
	GetFolder(ctx context.Context, id uuid.UUID, requestingUserID int) (domain.Folder, error)
}

func NewTasksService(tasksRepository TasksRepository, usersChecker UsersChecker, foldersChecker FoldersChecker) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
		usersChecker:    usersChecker,
		foldersChecker:  foldersChecker,
	}
}
