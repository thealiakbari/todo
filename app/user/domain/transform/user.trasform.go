package transform

import (
	"github.com/google/uuid"
	"github.com/thealiakbari/hichapp/app/user/domain/dto"
	"github.com/thealiakbari/hichapp/internal/user/domain/entity"
)

func CreateUserRequestToEntity(in dto.CreateUserRequest) entity.User {
	out := entity.User{
		Email: in.Email,
	}

	return out
}

func UpdateUserRequestToEntity(in dto.UpdateUserRequest, id string) (out entity.User, err error) {
	out = entity.User{
		Email: in.Email,
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return out, err
	}

	out.Id = idUUID
	return out, nil
}

func UserEntityToUserDto(in entity.User) dto.User {
	return dto.User{
		Id:        in.Id,
		Email:     in.Email,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func UsersEntityToUsersDto(in []entity.User) []dto.User {
	items := make([]dto.User, 0, len(in))
	for _, v := range in {
		items = append(items, UserEntityToUserDto(v))
	}

	return items
}
