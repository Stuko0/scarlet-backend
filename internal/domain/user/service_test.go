package user

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	userv1 "github.com/Stuko0/scarlet-backend/gen/proto/user/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctx        context.Context
	testSetup  *testUserService
	mockRepo   *MockUserRepository
	mockJWT    *MockJWTManager
}

type JWTManager interface {
	Generate(user *User) (string, error)
}

type MockUserRepository struct {
	mock.Mock
}

type MockJWTManager struct {
	mock.Mock
}

type testUserService struct {
	repo       *MockUserRepository
	jwtManager *MockJWTManager
	service    *UserService
}

func (s *UserServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockRepo = new(MockUserRepository)
	s.mockJWT = new(MockJWTManager)
	s.testSetup = &testUserService{
		repo:       s.mockRepo,
		jwtManager: s.mockJWT,
		service:    NewUserService(s.mockRepo, s.mockJWT),
	}
}

func (s *UserServiceTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
	s.mockJWT.AssertExpectations(s.T())
}

// Run the suite
func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (m *MockJWTManager) Generate(user *User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) CreateUserByEmail(ctx context.Context, user *User) (*User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) CreateUserByPhone(ctx context.Context, user *User) (*User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (s *UserServiceTestSuite) TestCreateUserByEmail_ValidRequest_Success() {
	request := &connect.Request[userv1.CreateUserByEmailRequest]{
		Msg: &userv1.CreateUserByEmailRequest{
			Name:     "John",
			Lastname: "Doe",
			Email:    "john@example.com",
			Password: "dont_care",
			Origin:   "web",
		},
	}

	expectedUser := &User{
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

	s.mockRepo.On("CreateUserByEmail", s.ctx, mock.MatchedBy(func(u *User) bool {
		return u.Name == request.Msg.Name &&
			u.Lastname == request.Msg.Lastname &&
			u.Email == request.Msg.Email &&
			u.Origin == request.Msg.Origin &&
			u.Password != request.Msg.Password
	})).Return(expectedUser, nil)

	response, err := s.testSetup.service.CreateUserByEmail(s.ctx, request)

	s.NoError(err)
	s.NotNil(response)
	s.Equal(expectedUser.UserId, response.Msg.User.UserId)
	s.Equal(expectedUser.Name, response.Msg.User.Name)
	s.Equal(expectedUser.Email, response.Msg.User.Email)
}

func (s *UserServiceTestSuite) TestCreateUserByEmail_EmptyCredentials_ReturnsError() {
	emptyEmailRequest := &connect.Request[userv1.CreateUserByEmailRequest]{
		Msg: &userv1.CreateUserByEmailRequest{
			Name:     "John",
			Lastname: "Doe",
			Email:    "",
			Password: "dont_care",
			Origin:   "web",
		},
	}

	response, err := s.testSetup.service.CreateUserByEmail(s.ctx, emptyEmailRequest)

	s.Nil(response)
	s.Error(err)
	connectErr := new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
	s.Contains(connectErr.Message(), "email and password are required")

	emptyPasswordRequest := &connect.Request[userv1.CreateUserByEmailRequest]{
		Msg: &userv1.CreateUserByEmailRequest{
			Name:     "John",
			Lastname: "Doe",
			Email:    "john@example.com",
			Password: "",
			Origin:   "web",
		},
	}
	response, err = s.testSetup.service.CreateUserByEmail(s.ctx, emptyPasswordRequest)

	s.Nil(response)
	s.Error(err)
	connectErr = new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
	s.Contains(connectErr.Message(), "email and password are required")
}

func (s *UserServiceTestSuite) TestCreateUserByEmail_ShortPassword_ReturnsError() {
	request := &connect.Request[userv1.CreateUserByEmailRequest]{
		Msg: &userv1.CreateUserByEmailRequest{
			Name:     "John",
			Lastname: "Doe",
			Email:    "john@example.com",
			Password: "short",
			Origin:   "web",
		},
	}

	response, err := s.testSetup.service.CreateUserByEmail(s.ctx, request)

	s.Nil(response)
	s.Error(err)
	connectErr := new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
	s.Contains(connectErr.Message(), "password must be 8+ characters")
}

func (s *UserServiceTestSuite) TestCreateUserByEmail_RepositoryError_ReturnsError() {
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
	s.mockRepo.On("CreateUserByEmail", s.ctx, mock.Anything).Return(nil, expectedError)

	response, err := s.testSetup.service.CreateUserByEmail(s.ctx, request)

	s.Nil(response)
	s.Error(err)
	connectErr := new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInternal, connectErr.Code())
}

// Converted Test Cases - LoginByEmail
func (s *UserServiceTestSuite) TestLoginByEmail_ValidCredentials_Success() {
	email := "john@example.com"
	password := "dont_care"

	request := &connect.Request[userv1.LoginByEmailRequest]{
		Msg: &userv1.LoginByEmailRequest{
			Email:    email,
			Password: password,
		},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := &User{
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

	s.mockRepo.On("GetUserByEmail", s.ctx, email).Return(mockUser, nil)
	s.mockJWT.On("Generate", mockUser).Return(mockToken, nil)

	response, err := s.testSetup.service.LoginByEmail(s.ctx, request)

	s.NoError(err)
	s.NotNil(response)
	s.Equal(mockToken, response.Msg.Token)
	s.Equal(mockUser.UserId, response.Msg.User.UserId)
	s.Equal(mockUser.Email, response.Msg.User.Email)
}

func (s *UserServiceTestSuite) TestLoginByEmail_EmptyCredentials_ReturnsError() {
	emptyEmailRequest := &connect.Request[userv1.LoginByEmailRequest]{
		Msg: &userv1.LoginByEmailRequest{
			Email:    "",
			Password: "dont_care",
		},
	}

	response, err := s.testSetup.service.LoginByEmail(s.ctx, emptyEmailRequest)

	s.Nil(response)
	s.Error(err)
	connectErr := new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
	s.Contains(connectErr.Message(), "email and password are required")

	emptyPasswordRequest := &connect.Request[userv1.LoginByEmailRequest]{
		Msg: &userv1.LoginByEmailRequest{
			Email:    "john@example.com",
			Password: "",
		},
	}

	response, err = s.testSetup.service.LoginByEmail(s.ctx, emptyPasswordRequest)

	s.Nil(response)
	s.Error(err)
	connectErr = new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
	s.Contains(connectErr.Message(), "email and password are required")
}

func (s *UserServiceTestSuite) TestLoginByEmail_UserNotFound_ReturnsError() {
	request := &connect.Request[userv1.LoginByEmailRequest]{
		Msg: &userv1.LoginByEmailRequest{
			Email:    "nonexistent@example.com",
			Password: "dont_care",
		},
	}

	s.mockRepo.On("GetUserByEmail", s.ctx, request.Msg.Email).
		Return(nil, errors.New("user not found"))

	response, err := s.testSetup.service.LoginByEmail(s.ctx, request)

	s.Nil(response)
	s.Error(err)
	connectErr := new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeNotFound, connectErr.Code())
	s.Contains(connectErr.Message(), "user not found")
}

func (s *UserServiceTestSuite) TestLoginByEmail_InvalidPassword_ReturnsError() {
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

	mockUser := &User{
		UserId:   1,
		Email:    email,
		Password: string(hashedPassword),
	}

	s.mockRepo.On("GetUserByEmail", s.ctx, email).Return(mockUser, nil)

	response, err := s.testSetup.service.LoginByEmail(s.ctx, request)

	s.Nil(response)
	s.Error(err)
	connectErr := new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeUnauthenticated, connectErr.Code())
	s.Contains(connectErr.Message(), "invalid credentials")
}

func (s *UserServiceTestSuite) TestLoginByEmail_JWTGenerationError_ReturnsError() {
	email := "john@example.com"
	password := "dont_care"

	request := &connect.Request[userv1.LoginByEmailRequest]{
		Msg: &userv1.LoginByEmailRequest{
			Email:    email,
			Password: password,
		},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := &User{
		UserId:   1,
		Email:    email,
		Password: string(hashedPassword),
	}

	jwtError := errors.New("failed to generate token")

	s.mockRepo.On("GetUserByEmail", s.ctx, email).Return(mockUser, nil)
	s.mockJWT.On("Generate", mockUser).Return("", jwtError)

	response, err := s.testSetup.service.LoginByEmail(s.ctx, request)

	s.Nil(response)
	s.Error(err)
	connectErr := new(connect.Error)
	s.True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInternal, connectErr.Code())
	s.Contains(connectErr.Message(), "failed to generate token")
}