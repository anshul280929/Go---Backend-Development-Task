package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/user/user-management-api/db/sqlc"
	"github.com/user/user-management-api/internal/logger"
	"github.com/user/user-management-api/internal/models"
	"github.com/user/user-management-api/internal/repository"
)

// Common sentinel errors for the handler layer to inspect.
var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("user not found")
	ErrInternal   = errors.New("internal server error")
)

// UserService contains business logic for user operations.
type UserService struct {
	repo     repository.UserRepository
	validate *validator.Validate
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo:     repo,
		validate: validator.New(),
	}
}

// CreateUser validates input, persists a new user, and returns the response.
func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	log := logger.Get()

	if err := s.validate.Struct(req); err != nil {
		log.Warn("validation failed for CreateUser", zap.Error(err))
		return nil, errors.Join(ErrValidation, err)
	}

	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		log.Warn("invalid date format", zap.String("dob", req.DOB), zap.Error(err))
		return nil, errors.Join(ErrValidation, errors.New("invalid date format, expected YYYY-MM-DD"))
	}

	user, err := s.repo.CreateUser(ctx, sqlc.CreateUserParams{
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		log.Error("failed to create user", zap.Error(err))
		return nil, errors.Join(ErrInternal, err)
	}

	log.Info("user created", zap.Int32("id", user.ID), zap.String("name", user.Name))

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
	}, nil
}

// GetUser retrieves a user by ID and calculates their age.
func (s *UserService) GetUser(ctx context.Context, id int32) (*models.UserResponse, error) {
	log := logger.Get()

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		log.Error("failed to get user", zap.Int32("id", id), zap.Error(err))
		return nil, errors.Join(ErrInternal, err)
	}

	age := models.CalculateAge(user.Dob)

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
		Age:  &age,
	}, nil
}

// ListUsers returns a paginated list of users with calculated ages.
func (s *UserService) ListUsers(ctx context.Context, params models.PaginationParams) (*models.PaginatedResponse, error) {
	log := logger.Get()

	params.Normalize()

	total, err := s.repo.CountUsers(ctx)
	if err != nil {
		log.Error("failed to count users", zap.Error(err))
		return nil, errors.Join(ErrInternal, err)
	}

	users, err := s.repo.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(params.PageSize),
		Offset: int32(params.Offset()),
	})
	if err != nil {
		log.Error("failed to list users", zap.Error(err))
		return nil, errors.Join(ErrInternal, err)
	}

	var responses []models.UserResponse
	for _, u := range users {
		age := models.CalculateAge(u.Dob)
		responses = append(responses, models.UserResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.Dob.Format("2006-01-02"),
			Age:  &age,
		})
	}

	// Return empty array instead of null
	if responses == nil {
		responses = []models.UserResponse{}
	}

	result := models.NewPaginatedResponse(responses, params, total)
	return &result, nil
}

// UpdateUser validates input, updates the user, and returns the response.
func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
	log := logger.Get()

	if err := s.validate.Struct(req); err != nil {
		log.Warn("validation failed for UpdateUser", zap.Error(err))
		return nil, errors.Join(ErrValidation, err)
	}

	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		log.Warn("invalid date format", zap.String("dob", req.DOB), zap.Error(err))
		return nil, errors.Join(ErrValidation, errors.New("invalid date format, expected YYYY-MM-DD"))
	}

	user, err := s.repo.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		log.Error("failed to update user", zap.Int32("id", id), zap.Error(err))
		return nil, errors.Join(ErrInternal, err)
	}

	log.Info("user updated", zap.Int32("id", user.ID), zap.String("name", user.Name))

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
	}, nil
}

// DeleteUser removes a user by ID.
func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	log := logger.Get()

	// Check existence first so we can return 404 properly.
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		log.Error("failed to check user existence", zap.Int32("id", id), zap.Error(err))
		return errors.Join(ErrInternal, err)
	}

	if err := s.repo.DeleteUser(ctx, id); err != nil {
		log.Error("failed to delete user", zap.Int32("id", id), zap.Error(err))
		return errors.Join(ErrInternal, err)
	}

	log.Info("user deleted", zap.Int32("id", id))
	return nil
}
