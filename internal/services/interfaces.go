package services

import (
	"context"

	"github.com/Stuko0/scarlet-backend/internal/models"
)

type UserRepositoryInterface interface{
	CreateUserByEmail(ctx context.Context, user *models.User)(*models.User, error)
	CreateUserByPhone(ctx context.Context, user *models.User)(*models.User, error)
	UpdateUser(ctx context.Context, user *models.User)(*models.User, error)
	GetUserByEmail(ctx context.Context, email string)(*models.User, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
}

type JWTManagerInterface interface{
	Generate(useer *models.User)(string, error)
}