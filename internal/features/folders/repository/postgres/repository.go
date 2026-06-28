package folders_postgres_repository

import (
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

type FoldersRepository struct {
	pool core_postgres_pool.Pool
}

func NewFoldersRepository(pool core_postgres_pool.Pool) *FoldersRepository {
	return &FoldersRepository{pool: pool}
}
