package transform

import (
	"github.com/google/uuid"
	"github.com/thealiakbari/todoapp/internal/application/todo/domain/dto"
	"github.com/thealiakbari/todoapp/internal/domain/todo/entity"
)

func CreateTodoItemRequestToEntity(in dto.CreateTodoItemRequest) entity.TodoItem {
	out := entity.TodoItem{
		Description: in.Description,
		DueDate:     in.DueDate,
	}

	return out
}

func UpdateTodoItemRequestToEntity(in dto.UpdateTodoItemRequest, id string) (out entity.TodoItem, err error) {
	out = entity.TodoItem{
		Description: in.Description,
		DueDate:     in.DueDate,
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return out, err
	}

	out.Id = idUUID
	return out, nil
}

func TodoItemEntityToTodoItemDto(in entity.TodoItem) dto.TodoItem {
	return dto.TodoItem{
		Id:          in.Id,
		Description: in.Description,
		DueDate:     in.DueDate,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

func TodoItemsEntityToTodoItemsDto(in []entity.TodoItem) []dto.TodoItem {
	items := make([]dto.TodoItem, 0, len(in))
	for _, v := range in {
		items = append(items, TodoItemEntityToTodoItemDto(v))
	}

	return items
}
