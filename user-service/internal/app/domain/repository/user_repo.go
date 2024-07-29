package repository

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/zer0day88/cme/user-service/internal/app/domain/entities"
)

type UserRepository struct {
	db *gocql.Session
}

func NewUserRepository(db *gocql.Session) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Get(ctx context.Context, users *[]entities.User) error {
	scanner := r.db.Query("select id,username,password,created_at from users").WithContext(ctx).Iter().Scanner()
	for scanner.Next() {
		var u entities.User
		err := scanner.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {

			return err
		}
		*users = append(*users, u)
	}
	if err := scanner.Err(); err != nil {

		return err
	}

	return nil
}

func (r *UserRepository) FindOneByUsername(ctx context.Context, username string) (*entities.User, error) {
	scanner := r.db.Query("select id,username,password,created_at from users where username=? ", username).
		WithContext(ctx).
		Consistency(gocql.One)

	var user entities.User

	err := scanner.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Insert(ctx context.Context, user entities.User) error {
	err := r.db.Query("insert into users(id,username,password,created_at) values (?,?,?,?)",
		user.ID, user.Username, user.Password, user.CreatedAt).
		WithContext(ctx).Exec()

	if err != nil {
		return err
	}

	return nil
}
