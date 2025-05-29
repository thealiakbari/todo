package repository

import (
	"context"

	"github.com/thealiakbari/hichapp/internal/poll/domain/entity"
	"github.com/thealiakbari/hichapp/pkg/common/db"
)

type PollRepository interface {
	Create(ctx context.Context, in entity.Poll) (res entity.Poll, err error)
	Update(ctx context.Context, in entity.Poll) (err error)
	FindByIds(ctx context.Context, ids []string) (res []entity.Poll, err error)
	FindByIdOrEmpty(ctx context.Context, id string) (res entity.Poll, err error)
	Purge(ctx context.Context, id string) (err error)
	Delete(ctx context.Context, id string) (err error)
	FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.Poll, err error)
	FilterCount(ctx context.Context, query []any) (res int64, err error)
}

type pollConfig struct {
	db db.DBWrapper
}

func NewPollRepository(db db.DBWrapper) PollRepository {
	return pollConfig{
		db: db,
	}
}

func (u pollConfig) Create(ctx context.Context, in entity.Poll) (res entity.Poll, err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return entity.Poll{}, err
	}

	return in, nil
}

func (u pollConfig) Update(ctx context.Context, in entity.Poll) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return err
	}

	return nil
}

func (u pollConfig) FindByIdOrEmpty(ctx context.Context, id string) (res entity.Poll, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order("created_at desc").Find(&res, "id = ?", id).Limit(1).Error
	if err != nil {
		return entity.Poll{}, err
	}

	return res, nil
}

func (u pollConfig) FindByIds(ctx context.Context, ids []string) (res []entity.Poll, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Find(&res, "id IN (?)", ids).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u pollConfig) Purge(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Exec("DELETE FROM polls WHERE id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u pollConfig) Delete(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Delete(&entity.Poll{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u pollConfig) FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.Poll, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order(order).
		Limit(limit).
		Offset(offset).
		Find(&res, query...).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u pollConfig) FilterCount(ctx context.Context, query []any) (res int64, err error) {
	countQuery := db.GormConnection(ctx, u.db.DB).Model(&entity.Poll{})
	if len(query) > 1 {
		countQuery = countQuery.Where(query[0], query[1:]...)
	} else if len(query) == 1 {
		countQuery = countQuery.Where(query[0])
	}

	err = countQuery.Count(&res).Error
	if err != nil {
		return 0, err
	}

	return res, nil
}
