package pg

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thealiakbari/todoapp/internal/domain/todo/entity"
	"github.com/thealiakbari/todoapp/pkg/common/config"
	"github.com/thealiakbari/todoapp/pkg/common/db"
)

func setupTestDB(t *testing.T) db.DBWrapper {
	conf := config.LoadConfig("../../../../../config/todoapp.yml")
	gormDB, err := db.NewPostgresConn(context.Background(), conf.DB.Postgres)
	assert.NoError(t, err)
	dbw := db.NewDBWrapper(gormDB)

	return dbw
}

func TestTodoItemRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	testDB := setupTestDB(t)
	repo := NewTodoItemRepository(testDB)

	// Create
	item := entity.TodoItem{
		Description: "Test Task",
		DueDate:     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
	created, err := repo.Create(ctx, item)
	assert.NoError(t, err)
	assert.Equal(t, item.Description, created.Description)

	// FindByIdOrEmpty
	found, err := repo.FindByIdOrEmpty(ctx, created.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, created.Id, found.Id)

	// Update
	created.Description = "Updated Task"
	err = repo.Update(ctx, created)
	assert.NoError(t, err)

	updated, err := repo.FindByIdOrEmpty(ctx, created.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", updated.Description)

	// FindByIds
	list, err := repo.FindByIds(ctx, []string{created.Id.String()})
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// FilterFind
	results, err := repo.FilterFind(ctx, []any{"description LIKE ?", "%Task%"}, "created_at desc", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	// FilterCount
	count, err := repo.FilterCount(ctx, []any{"description LIKE ?", "%Task%"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Delete
	err = repo.Delete(ctx, created.Id.String())
	assert.NoError(t, err)

	_, err = repo.FindByIdOrEmpty(ctx, created.Id.String())
	assert.NoError(t, err) // should return empty entity, not fail

	// Purge
	// Re-create and then purge
	item2 := entity.TodoItem{
		Description: "Temp Task",
		DueDate:     time.Now().Format(time.RFC3339),
	}
	created2, _ := repo.Create(ctx, item2)
	err = repo.Purge(ctx, created2.Id.String())
	assert.NoError(t, err)
}
