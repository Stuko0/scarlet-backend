package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Stuko0/scarlet-backend/internal/database"
	"github.com/jackc/pgx/v5"
)

type UserRepositoryInterface interface{
	CreateUserByEmail(ctx context.Context, user *User)(*User, error)
	CreateUserByPhone(ctx context.Context, user *User)(*User, error)
	UpdateUser(ctx context.Context, user *User)(*User, error)
	GetUserByEmail(ctx context.Context, email string)(*User, error)
	GetUserByID(ctx context.Context, userID int64) (*User, error)
}

type JWTManagerInterface interface{
	Generate(useer *User)(string, error)
}

type UserRepository struct {
	db *database.Postgres
}

func NewUserRepository(db *database.Postgres) *UserRepository {
	return &UserRepository{db: db,}
}

func (r *UserRepository) CreateUserByEmail(ctx context.Context, user *User) (*User,error) {
	query := `INSERT INTO scarlet.users (name, lastname, email, password, origin) 
	VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
	var createdUser User
	err := r.db.Pool.QueryRow(ctx, query,
		user.Name,
		user.Lastname,
		user.Email,
		user.Password,
		user.Origin,).Scan(&createdUser.UserId)
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

func (r *UserRepository) CreateUserByPhone(ctx context.Context, user *User)(*User, error){
	query:= `INSERT INTO scarlet.users (name, lastname, email, phone, origin) VALUES ($1,$2,$3,$4,$5) RETURNING user_id`
	var createdUser User
	err:=r.db.Pool.QueryRow(ctx, query,
		user.Name,
		user.Lastname,
		user.Email,
		user.Phone,
		user.Origin,).Scan(&createdUser.UserId)
	if err!=nil{return nil, err}
	return &createdUser, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *User)(*User, error){
	query:=`UPDATE scarlet.users SET name=$1,lastname=$2,email=$3,phone=$4,password=$5,role=$6,image=$7, updated_at=NOW() WHERE user_id=$8`
	_,err:=r.db.Pool.Exec(ctx,query,
		user.Name,
		user.Lastname,
		user.Email,
		user.Phone,
		user.Password,
		user.Role,
		user.Image,
	)
	if err!= nil{
		return nil, fmt.Errorf("failed to update user %v", err)
	}
	return r.GetUserByID(ctx, user.UserId)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string)(*User, error){
	query := `SELECT user_id, name, lastname, email, password, role FROM scarlet.users WHERE email=$1`
	var user User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.UserId,
		&user.Name,
		&user.Lastname,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err !=nil{return nil, err}
	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int64) (*User, error){
	query:=`SELECT user_id, name, lastname, email, role FROM scarlet.users WHERE user_id=$1`
	var user User
	err:=r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.UserId,
		&user.Name,
		&user.Email,
		&user.Role,
	)
	if err!= nil{
		if errors.Is(err, pgx.ErrNoRows){
			return nil, ErrUserNotFound
		}
	}
	return &user, nil
}

var ErrUserNotFound =errors.New("user not found")