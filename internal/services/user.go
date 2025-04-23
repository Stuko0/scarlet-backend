package services

import (
	"context"
	"errors"
	"connectrpc.com/connect"
	userv1 "github.com/Stuko0/scarlet-backend/gen/proto/user/v1"
	// "github.com/Stuko0/scarlet-backend/internal/auth"
	"github.com/Stuko0/scarlet-backend/internal/models"
	"github.com/Stuko0/scarlet-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{
	repo UserRepositoryInterface
	jwtManager JWTManagerInterface
}

func userToProto(user *models.User) *userv1.User {
	return &userv1.User{
		UserId:    user.UserId,
		Name:      user.Name,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Phone:     user.Phone,
		Origin:    user.Origin,
		Active:    user.Active,
		Role:      user.Role,
		Image:     user.Image,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserService(repo UserRepositoryInterface, jwtManager JWTManagerInterface) *UserService {
	return &UserService{
		repo: repo,
		jwtManager: jwtManager,
	}
}
func (s *UserService) CreateUserByEmail(ctx context.Context, req *connect.Request[userv1.CreateUserByEmailRequest],) (*connect.Response[userv1.UserResponse], error) {
	if req.Msg.Email=="" || req.Msg.Password==""{
		return nil, connect.NewError(connect.CodeInvalidArgument, connect.NewError(connect.CodeInvalidArgument, errors.New("email and password are required")))
	}

	if len(req.Msg.Password)<8{
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password must be 8+ characters"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Msg.Password),
		bcrypt.DefaultCost,
	)

	if err != nil{return nil, connect.NewError(connect.CodeInternal, err)}

	userModel := &models.User{
		Name:     req.Msg.Name,
		Lastname: req.Msg.Lastname,
		Email:    req.Msg.Email,
		Password: string(hashedPassword),
		Origin:   req.Msg.Origin,
	}

	createdUser, err := s.repo.CreateUserByEmail(ctx, userModel)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res:= connect.NewResponse(&userv1.UserResponse{
		User: userToProto(createdUser),
	})
	return res,nil
}

func (s *UserService) CreateUserByPhone(ctx context.Context, req *connect.Request[userv1.CreateUserByPhoneRequest],)(*connect.Response[userv1.UserResponse], error){
	if req.Msg.Phone==""{
		return nil, connect.NewError(connect.CodeInvalidArgument, connect.NewError(connect.CodeInvalidArgument, errors.New("phone is required")))
	}
	userModel:= &models.User{
		Name: req.Msg.Name,
		Lastname: req.Msg.Lastname,
		Email: req.Msg.Email,
		Phone: req.Msg.Phone,
		Origin: req.Msg.Origin,
	}
	createdUser, err:= s.repo.CreateUserByPhone(ctx, userModel)
	if err!=nil{return nil, connect.NewError(connect.CodeInternal, err)}
	res:=connect.NewResponse(&userv1.UserResponse{
		User: userToProto(createdUser),
	})
	return res,nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *connect.Request[userv1.UpdateUserRequest],) (*connect.Response[userv1.UserResponse], error){
	if req.Msg.UserId==0{
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user ID is required"))
	}

	update:=&models.User{
		UserId: req.Msg.UserId,
		Name: req.Msg.Name,
		Lastname: req.Msg.Lastname,
		Email: req.Msg.Email,
		Phone: req.Msg.Phone,
		Password: req.Msg.Password,
		Role: req.Msg.Role,
		Image: req.Msg.Image,
	}

	updateUser,err:=s.repo.UpdateUser(ctx, update)
	if err!=nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
		}
		return nil, connect.NewError(connect.CodeInternal,err)
	}

	return connect.NewResponse(&userv1.UserResponse{
		User: userToProto(updateUser),
	}), nil
}

func (s *UserService) LoginByEmail (ctx context.Context, req *connect.Request[userv1.LoginByEmailRequest],)(*connect.Response[userv1.LoginByEmailResponse], error){
	if req.Msg.Email=="" || req.Msg.Password==""{
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email and password are required"))
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Msg.Email)
	if err != nil{
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	if err:= bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Msg.Password),);
	err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil{
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate token"))
	}

	res:=connect.NewResponse(&userv1.LoginByEmailResponse{
		Token: token,
		User: userToProto(user),
	})
	return res, nil
}

func (s *UserService)GetUser(ctx context.Context, req *connect.Request[userv1.GetUserRequest],)(*connect.Response[userv1.UserResponse], error){
	userID:=req.Msg.GetUserId()
	if userID ==0{
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user ID is required"))
	}

	user, err:= s.repo.GetUserByID(ctx, userID)
	if err!=nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&userv1.UserResponse{
		User: userToProto(user),
	}), nil
}

func (s *UserService) DeleteUser (ctx context.Context, req *connect.Request[userv1.DeleteUserRequest],) (*connect.Response[userv1.DeleteUserResponse], error){
	return connect.NewResponse(&userv1.DeleteUserResponse{Success: true}), nil
}