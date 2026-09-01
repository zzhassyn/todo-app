package tasks_postgres_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func TestCreateTask(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// 1. Setup prerequisite data (user)
	userID := 1
	setupUserSQL := `INSERT INTO todoapp.users (id, full_name, email, password_hash) VALUES ($1, 'Test User', 'test@example.com', 'hash')`
	_, err := repo.pool.Exec(ctx, setupUserSQL, userID)
	require.NoError(t, err)

	// 2. Setup a tag (just to test tag association, even though tags are normally managed by get_tags/etc)
	setupTagSQL := `INSERT INTO todoapp.tags (name) VALUES ('work'), ('urgent')`
	_, err = repo.pool.Exec(ctx, setupTagSQL)
	require.NoError(t, err)

	// 3. Test Create Task without tags
	t.Run("success without tags", func(t *testing.T) {
		desc := "Description"
		task := domain.NewTaskUninitialized("Test Task 1", &desc, userID, domain.PriorityMedium, nil, nil)
		created, err := repo.CreateTask(ctx, task, nil)
		assert.NoError(t, err)
		assert.Equal(t, "Test Task 1", created.Title)
		assert.Equal(t, "Description", *created.Description)
		assert.Equal(t, domain.PriorityMedium, created.Priority)
		assert.NotZero(t, created.ID)
		assert.NotZero(t, created.Position)
	})

	// 4. Test Create Task with tags
	t.Run("success with tags", func(t *testing.T) {
		task := domain.NewTaskUninitialized("Test Task 2", nil, userID, domain.PriorityHigh, nil, nil)
		created, err := repo.CreateTask(ctx, task, []string{"work", "urgent"})
		assert.NoError(t, err)
		assert.Equal(t, "Test Task 2", created.Title)
		assert.Nil(t, created.Description)
		assert.Len(t, created.Tags, 2)
	})
	
	// 5. Test invalid user ID
	t.Run("foreign key violation", func(t *testing.T) {
		task := domain.NewTaskUninitialized("Test Task 3", nil, 999, domain.PriorityLow, nil, nil)
		_, err := repo.CreateTask(ctx, task, nil)
		assert.Error(t, err)
	})
}
