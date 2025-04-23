package services_test

import (
	"context"
	"testing"
	"errors"

	userv1 "github.com/Stuko0/scarlet-backend/gen/proto/user/v1"
	"github.com/Stuko0/scarlet-backend/internal/models"
	"github.com/Stuko0/scarlet-backend/internal/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"connectrpc.com/connect"

)

type UserRepository interface {
	CreateUserByEmail(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type JWTManager interface {
	Generate(user *models.User) (string, error)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUserByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUserByPhone(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) Generate(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

type testUserService struct {
	repo       *MockUserRepository
	jwtManager *MockJWTManager
	service    *services.UserService
}

func setupTestService() *testUserService {
	mockRepo := &MockUserRepository{}
	mockJWT := &MockJWTManager{}
	
	service := services.NewUserService(mockRepo, mockJWT)
	
	return &testUserService{
		repo:       mockRepo,
		jwtManager: mockJWT,
		service:    service,
	}
}

func TestCreateUserByEmail(t *testing.T) {
	t.Run("ValidRequest_Success", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		request := &connect.Request[userv1.CreateUserByEmailRequest]{
			Msg: &userv1.CreateUserByEmailRequest{
				Name:     "John",
				Lastname: "Doe",
				Email:    "john@example.com",
				Password: "dont_care",
				Origin:   "web",
			},
		}
		
		expectedUser := &models.User{
			UserId:    1,
			Name:      "John",
			Lastname:  "Doe",
			Email:     "john@example.com",
			Password:  "dont_care",
			Origin:    "web",
			Active:    true,
			Role:      "user",
			CreatedAt: "2025-04-23T00:00:00Z",
			UpdatedAt: "2025-04-23T00:00:00Z",
		}
		
		testSetup.repo.On("CreateUserByEmail", ctx, mock.MatchedBy(func(u *models.User) bool {
			return u.Name == request.Msg.Name &&
				u.Lastname == request.Msg.Lastname &&
				u.Email == request.Msg.Email &&
				u.Origin == request.Msg.Origin &&
				u.Password != request.Msg.Password
		})).Return(expectedUser, nil)
		
		response, err := testSetup.service.CreateUserByEmail(ctx, request)
		
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedUser.UserId, response.Msg.User.UserId)
		assert.Equal(t, expectedUser.Name, response.Msg.User.Name)
		assert.Equal(t, expectedUser.Email, response.Msg.User.Email)
		
		testSetup.repo.AssertExpectations(t)
	})
	
	t.Run("EmptyCredentials_ReturnsError", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		emptyEmailRequest := &connect.Request[userv1.CreateUserByEmailRequest]{
			Msg: &userv1.CreateUserByEmailRequest{
				Name:     "John",
				Lastname: "Doe",
				Email:    "",
				Password: "dont_care",
				Origin:   "web",
			},
		}
		
		response, err := testSetup.service.CreateUserByEmail(ctx, emptyEmailRequest)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok := err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
		assert.Contains(t, connectErr.Message(), "email and password are required")
		
		emptyPasswordRequest := &connect.Request[userv1.CreateUserByEmailRequest]{
			Msg: &userv1.CreateUserByEmailRequest{
				Name:     "John",
				Lastname: "Doe",
				Email:    "john@example.com",
				Password: "",
				Origin:   "web",
			},
		}
		response, err = testSetup.service.CreateUserByEmail(ctx, emptyPasswordRequest)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok = err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
		assert.Contains(t, connectErr.Message(), "email and password are required")
	})
	
	t.Run("ShortPassword_ReturnsError", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		request := &connect.Request[userv1.CreateUserByEmailRequest]{
			Msg: &userv1.CreateUserByEmailRequest{
				Name:     "John",
				Lastname: "Doe",
				Email:    "john@example.com",
				Password: "short",
				Origin:   "web",
			},
		}
		
		response, err := testSetup.service.CreateUserByEmail(ctx, request)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok := err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
		assert.Contains(t, connectErr.Message(), "password must be 8+ characters")
	})
	
	t.Run("RepositoryError_ReturnsError", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		request := &connect.Request[userv1.CreateUserByEmailRequest]{
			Msg: &userv1.CreateUserByEmailRequest{
				Name:     "John",
				Lastname: "Doe",
				Email:    "john@example.com",
				Password: "dont_care",
				Origin:   "web",
			},
		}
		
		expectedError := errors.New("database error")
		testSetup.repo.On("CreateUserByEmail", ctx, mock.Anything).Return(nil, expectedError)
		
		response, err := testSetup.service.CreateUserByEmail(ctx, request)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok := err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeInternal, connectErr.Code())
		
		testSetup.repo.AssertExpectations(t)
	})
}

func TestLoginByEmail(t *testing.T) {
	t.Run("ValidCredentials_Success", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		email := "john@example.com"
		password := "dont_care"
		
		request := &connect.Request[userv1.LoginByEmailRequest]{
			Msg: &userv1.LoginByEmailRequest{
				Email:    email,
				Password: password,
			},
		}
		
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		mockUser := &models.User{
			UserId:    1,
			Name:      "John",
			Lastname:  "Doe",
			Email:     email,
			Password:  string(hashedPassword),
			Origin:    "web",
			Active:    true,
			Role:      "user",
			CreatedAt: "2025-04-23T00:00:00Z",
			UpdatedAt: "2025-04-23T00:00:00Z",
		}
		
		mockToken := "jwt.mock.token"
		
		testSetup.repo.On("GetUserByEmail", ctx, email).Return(mockUser, nil)
		testSetup.jwtManager.On("Generate", mockUser).Return(mockToken, nil)
		
		response, err := testSetup.service.LoginByEmail(ctx, request)
		
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, mockToken, response.Msg.Token)
		assert.Equal(t, mockUser.UserId, response.Msg.User.UserId)
		assert.Equal(t, mockUser.Email, response.Msg.User.Email)
		
		testSetup.repo.AssertExpectations(t)
		testSetup.jwtManager.AssertExpectations(t)
	})
	
	t.Run("EmptyCredentials_ReturnsError", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		emptyEmailRequest := &connect.Request[userv1.LoginByEmailRequest]{
			Msg: &userv1.LoginByEmailRequest{
				Email:    "",
				Password: "dont_care",
			},
		}
		
		response, err := testSetup.service.LoginByEmail(ctx, emptyEmailRequest)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok := err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
		assert.Contains(t, connectErr.Message(), "email and password are required")
		
		emptyPasswordRequest := &connect.Request[userv1.LoginByEmailRequest]{
			Msg: &userv1.LoginByEmailRequest{
				Email:    "john@example.com",
				Password: "",
			},
		}
		
		response, err = testSetup.service.LoginByEmail(ctx, emptyPasswordRequest)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok = err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
		assert.Contains(t, connectErr.Message(), "email and password are required")
	})
	
	t.Run("UserNotFound_ReturnsError", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		request := &connect.Request[userv1.LoginByEmailRequest]{
			Msg: &userv1.LoginByEmailRequest{
				Email:    "nonexistent@example.com",
				Password: "dont_care",
			},
		}
		
		testSetup.repo.On("GetUserByEmail", ctx, request.Msg.Email).
			Return(nil, errors.New("user not found"))
		
		response, err := testSetup.service.LoginByEmail(ctx, request)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok := err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeNotFound, connectErr.Code())
		assert.Contains(t, connectErr.Message(), "user not found")
		
		testSetup.repo.AssertExpectations(t)
	})
	
	t.Run("InvalidPassword_ReturnsError", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		email := "john@example.com"
		correctPassword := "dont_care"
		wrongPassword := "wrongpassword"
		
		request := &connect.Request[userv1.LoginByEmailRequest]{
			Msg: &userv1.LoginByEmailRequest{
				Email:    email,
				Password: wrongPassword,
			},
		}
		
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
		
		mockUser := &models.User{
			UserId:    1,
			Email:     email,
			Password:  string(hashedPassword),
		}
		
		testSetup.repo.On("GetUserByEmail", ctx, email).Return(mockUser, nil)
		
		 
		response, err := testSetup.service.LoginByEmail(ctx, request)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok := err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeUnauthenticated, connectErr.Code())
		assert.Contains(t, connectErr.Message(), "invalid credentials")
		
		testSetup.repo.AssertExpectations(t)
	})
	
	t.Run("JWTGenerationError_ReturnsError", func(t *testing.T) {
		ctx := context.Background()
		testSetup := setupTestService()
		
		email := "john@example.com"
		password := "dont_care"
		
		request := &connect.Request[userv1.LoginByEmailRequest]{
			Msg: &userv1.LoginByEmailRequest{
				Email:    email,
				Password: password,
			},
		}
		
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		mockUser := &models.User{
			UserId:    1,
			Email:     email,
			Password:  string(hashedPassword),
		}
		
		jwtError := errors.New("failed to generate token")
		
		testSetup.repo.On("GetUserByEmail", ctx, email).Return(mockUser, nil)
		testSetup.jwtManager.On("Generate", mockUser).Return("", jwtError)
		
		 
		response, err := testSetup.service.LoginByEmail(ctx, request)
		
		assert.Nil(t, response)
		assert.Error(t, err)
		connectErr, ok := err.(*connect.Error)
		assert.True(t, ok)
		assert.Equal(t, connect.CodeInternal, connectErr.Code())
		assert.Contains(t, connectErr.Message(), "failed to generate token")
		
		testSetup.repo.AssertExpectations(t)
		testSetup.jwtManager.AssertExpectations(t)
	})
}