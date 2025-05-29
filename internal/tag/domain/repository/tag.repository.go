package repository

import (
	"context"

	"github.com/thealiakbari/hichapp/internal/tag/domain/entity"
	"github.com/thealiakbari/hichapp/pkg/common/db"
)

type TagRepository interface {
	Create(ctx context.Context, in entity.Tag) (res entity.Tag, err error)
	Update(ctx context.Context, in entity.Tag) (err error)
	FindByIds(ctx context.Context, ids []string) (res []entity.Tag, err error)
	FindByIdOrEmpty(ctx context.Context, id string) (res entity.Tag, err error)
	Purge(ctx context.Context, id string) (err error)
	Delete(ctx context.Context, id string) (err error)
	FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.Tag, err error)
	FilterCount(ctx context.Context, query []any) (res int64, err error)
}

type tagConfig struct {
	db db.DBWrapper
}

func NewTagRepository(db db.DBWrapper) TagRepository {
	return tagConfig{
		db: db,
	}
}

func (u tagConfig) Create(ctx context.Context, in entity.Tag) (res entity.Tag, err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return entity.Tag{}, err
	}

	return in, nil
}

func (u tagConfig) Update(ctx context.Context, in entity.Tag) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return err
	}

	return nil
}

func (u tagConfig) FindByIdOrEmpty(ctx context.Context, id string) (res entity.Tag, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order("created_at desc").Find(&res, "id = ?", id).Limit(1).Error
	if err != nil {
		return entity.Tag{}, err
	}

	return res, nil
}

func (u tagConfig) FindByIds(ctx context.Context, ids []string) (res []entity.Tag, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Find(&res, "id IN (?)", ids).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u tagConfig) Purge(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Exec("DELETE FROM tags WHERE id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u tagConfig) Delete(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Delete(&entity.Tag{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u tagConfig) FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.Tag, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order(order).
		Limit(limit).
		Offset(offset).
		Find(&res, query...).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u tagConfig) FilterCount(ctx context.Context, query []any) (res int64, err error) {
	countQuery := db.GormConnection(ctx, u.db.DB).Model(&entity.Tag{})
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
