package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
	"strings"
	"user-service/internal/model"
)

type UserRepository interface {
	GetUserById(ctx context.Context, userId int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdateUserById(ctx context.Context, user model.User) (model.User, error)
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserList(ctx context.Context, limit, offset int) ([]model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func (u *userRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	query := `INSERT INTO users (email, password, name, gender, phone, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	row := u.db.QueryRowxContext(ctx, query, user.Email, user.Password, user.Name, user.Gender, user.Phone, user.Role)
	var insertedUser model.User
	err := row.StructScan(&insertedUser)
	if err != nil {
		log.Printf("User repository, create user err: %v", err)
		return model.User{}, err
	}
	return insertedUser, nil
}

func (u *userRepository) GetUserList(ctx context.Context, limit, offset int) ([]model.User, error) {
	query := `SELECT * FROM users ORDER BY user_id LIMIT $1 OFFSET $2`
	rows, err := u.db.QueryxContext(ctx, query, limit, offset)
	if err != nil {
		log.Printf("User repository, get user list err: %v", err)
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.StructScan(&user)
		if err != nil {
			log.Printf("User repository, get user list err: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *userRepository) GetUserById(ctx context.Context, userId int) (model.User, error) {
	query := `SELECT * FROM users WHERE user_id = $1`
	row := u.db.QueryRowxContext(ctx, query, userId)
	var user model.User
	err := row.StructScan(&user)
	if err != nil {
		log.Printf("User repository, get user by id err: %v", err)
		return model.User{}, err
	}
	return user, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	row := u.db.QueryRowxContext(ctx, query, email)
	var user model.User
	err := row.StructScan(&user)
	if err != nil {
		log.Printf("User repository, get user by email err: %v", err)
		return model.User{}, err
	}
	return user, nil
}

func (u *userRepository) UpdateUserById(ctx context.Context, user model.User) (model.User, error) {
	var updateFields []string
	var values []interface{}
	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)
	cnt := 1
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Name == "UserId" {
			continue
		}
		if !value.IsZero() {
			updateFields = append(updateFields, fmt.Sprintf("%s = $%d", field.Tag.Get("db"), cnt))
			values = append(values, value.Interface())
			cnt += 1
		}
	}
	if len(updateFields) == 0 {
		return model.User{}, nil
	}
	values = append(values, user.UserId)
	query := fmt.Sprintf("UPDATE users SET %s WHERE user_id = $%d RETURNING *", strings.Join(updateFields, ", "), cnt)
	row := u.db.QueryRowxContext(ctx, query, values...)
	var updatedUser model.User
	err := row.StructScan(&updatedUser)
	if err != nil {
		log.Printf("User repository, update users err: %v", err)
		return model.User{}, err
	}
	return updatedUser, nil
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}
