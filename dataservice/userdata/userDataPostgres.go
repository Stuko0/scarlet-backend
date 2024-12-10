package userdata

import (
	"database/sql"
	_ "github.com/lib/pq"
	"scarlet_backend/model"
)

type UserDataPostgres struct {
	db *sql.DB
}

func NewUserDataPostgres(db *sql.DB) *UserDataPostgres {
	return &UserDataPostgres{db: db}
}

func (u *UserDataPostgres) GetUsers() ([]model.User, error) {
	rows, err := u.db.Query("select * from user_data.users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.Id, &user.Name, &user.Lastname, &user.Email, &user.Phone, &user.Psw, &user.Origin, &user.Active, &user.CreatedAt, &user.Rol); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserDataPostgres) FindByEmail(email string) (*model.User, error) {
	query := `select * from user_data.users where email=$1`
	row := u.db.QueryRow(query, email)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Lastname, &user.Email, &user.Phone, &user.Psw, &user.Origin, &user.Active, &user.CreatedAt, &user.Rol)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserDataPostgres) SaveByEmail(user *model.User) (*model.User, error) {
	query := `INSERT INTO user_data.users (userName, lastname, email, phone, psw, origin, active, created_at, rol) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING userId`
	err := u.db.QueryRow(query, user.Name, user.Lastname, user.Email, user.Phone, user.Psw, user.Origin, user.Active, user.CreatedAt, user.Rol).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
