package subtasks_postgres_repository

import (
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

type SubtasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewSubtasksRepository(pool core_postgres_pool.Pool) *SubtasksRepository {
	return &SubtasksRepository{pool: pool}
}
