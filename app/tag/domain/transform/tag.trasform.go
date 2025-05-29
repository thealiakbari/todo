package transform

import (
	"github.com/google/uuid"
	"github.com/thealiakbari/hichapp/app/tag/domain/dto"
	"github.com/thealiakbari/hichapp/internal/tag/domain/entity"
)

func CreateTagRequestToEntity(in dto.CreateTagRequest) entity.Tag {
	out := entity.Tag{
		Name: in.Email,
	}

	return out
}

func UpdateTagRequestToEntity(in dto.UpdateTagRequest, id string) (out entity.Tag, err error) {
	out = entity.Tag{
		Name: in.Email,
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return out, err
	}

	out.Id = idUUID
	return out, nil
}

func TagEntityToTagDto(in entity.Tag) dto.Tag {
	return dto.Tag{
		Id:        in.Id,
		Name:      in.Name,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func TagsEntityToTagsDto(in []entity.Tag) []dto.Tag {
	items := make([]dto.Tag, 0, len(in))
	for _, v := range in {
		items = append(items, TagEntityToTagDto(v))
	}

	return items
}
