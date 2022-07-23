package repository

import (
	"context"

	"truckx/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userImpl{db}
}

type userImpl struct {
	db *sqlx.DB
}

func (r userImpl) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE * FROM users WHERE id = $1", id)
	return err
}

func (r userImpl) Get(ctx context.Context, id int) (*models.User, error) {
	usr := models.User{}
	row := r.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE id = $1", id)
	err := row.StructScan(&usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func (r userImpl) List(ctx context.Context) ([]models.User, error) {
	users := []models.User{}
	row, err := r.db.QueryxContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	var user models.User
	for row.Next() {
		err = row.StructScan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
