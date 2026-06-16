package repository

import (
	"context"

	"github.com/user/user-management-api/db/sqlc"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByID(ctx context.Context, id int32) (sqlc.User, error)
	ListUsers(ctx context.Context, arg sqlc.ListUsersParams) ([]sqlc.User, error)
	CountUsers(ctx context.Context) (int64, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	DeleteUser(ctx context.Context, id int32) error
}

// userRepository implements UserRepository using SQLC Queries.
type userRepository struct {
	queries *sqlc.Queries
}

// NewUserRepository creates a new UserRepository backed by SQLC.
func NewUserRepository(queries *sqlc.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, arg)
}

func (r *userRepository) GetUserByID(ctx context.Context, id int32) (sqlc.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *userRepository) ListUsers(ctx context.Context, arg sqlc.ListUsersParams) ([]sqlc.User, error) {
	return r.queries.ListUsers(ctx, arg)
}

func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *userRepository) UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	return r.queries.UpdateUser(ctx, arg)
}

func (r *userRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}
