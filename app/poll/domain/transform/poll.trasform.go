package transform

import (
	"github.com/google/uuid"
	"github.com/thealiakbari/hichapp/app/poll/domain/dto"
	"github.com/thealiakbari/hichapp/internal/poll/domain/entity"
)

func CreatePollRequestToEntity(in dto.CreatePollRequest) entity.Poll {
	out := entity.Poll{
		Title: in.Title,
	}

	return out
}

func UpdatePollRequestToEntity(in dto.UpdatePollRequest, id string) (out entity.Poll, err error) {
	out = entity.Poll{
		Title: in.Title,
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return out, err
	}

	out.Id = idUUID
	return out, nil
}

func PollEntityToPollDto(in entity.Poll) dto.Poll {
	return dto.Poll{
		Id:        in.Id,
		Title:     in.Title,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func PollsEntityToPollsDto(in []entity.Poll) []dto.Poll {
	items := make([]dto.Poll, 0, len(in))
	for _, v := range in {
		items = append(items, PollEntityToPollDto(v))
	}

	return items
}
