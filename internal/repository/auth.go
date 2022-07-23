package repository

import (
	"context"

	"truckx/internal/models"

	"github.com/jmoiron/sqlx"
)

type AuthRepo interface {
	Login(ctx context.Context, email string) (*models.User, error)
	Signup(ctx context.Context, user models.User) (*models.User, error)
}

func NewAuthRepo(db *sqlx.DB) AuthRepo {
	return &authImpl{db}
}

type authImpl struct {
	db *sqlx.DB
}

func (r authImpl) Login(ctx context.Context, email string) (*models.User, error) {
	usr := models.User{}
	row := r.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE email = $1", email)
	err := row.StructScan(&usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func (r authImpl) Signup(ctx context.Context, user models.User) (*models.User, error) {
	row, err := r.db.QueryContext(ctx,
		"INSERT INTO Users (name, email, password, phone_no, role_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Name, user.Email, user.Password, user.PhoneNo, user.RoleId)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
