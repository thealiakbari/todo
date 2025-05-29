package repository

import (
	"context"

	"github.com/thealiakbari/hichapp/internal/user/domain/entity"
	"github.com/thealiakbari/hichapp/pkg/common/db"
)

type UserRepository interface {
	Create(ctx context.Context, in entity.User) (res entity.User, err error)
	Update(ctx context.Context, in entity.User) (err error)
	FindByIds(ctx context.Context, ids []string) (res []entity.User, err error)
	FindByIdOrEmpty(ctx context.Context, id string) (res entity.User, err error)
	Purge(ctx context.Context, id string) (err error)
	Delete(ctx context.Context, id string) (err error)
	FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.User, err error)
	FilterCount(ctx context.Context, query []any) (res int64, err error)
}

type userConfig struct {
	db db.DBWrapper
}

func NewUserRepository(db db.DBWrapper) UserRepository {
	return userConfig{
		db: db,
	}
}

func (u userConfig) Create(ctx context.Context, in entity.User) (res entity.User, err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return entity.User{}, err
	}

	return in, nil
}

func (u userConfig) Update(ctx context.Context, in entity.User) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return err
	}

	return nil
}

func (u userConfig) FindByIdOrEmpty(ctx context.Context, id string) (res entity.User, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order("created_at desc").Find(&res, "id = ?", id).Limit(1).Error
	if err != nil {
		return entity.User{}, err
	}

	return res, nil
}

func (u userConfig) FindByIds(ctx context.Context, ids []string) (res []entity.User, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Find(&res, "id IN (?)", ids).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u userConfig) Purge(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Exec("DELETE FROM users WHERE id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u userConfig) Delete(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Delete(&entity.User{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u userConfig) FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.User, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order(order).
		Limit(limit).
		Offset(offset).
		Find(&res, query...).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u userConfig) FilterCount(ctx context.Context, query []any) (res int64, err error) {
	countQuery := db.GormConnection(ctx, u.db.DB).Model(&entity.User{})
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
