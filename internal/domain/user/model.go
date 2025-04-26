package user

type User struct {
	UserId   int64 `json:"user_id"`
	Name string `db:"name"`
	Lastname string `db:"lastname"`
	Email string `db:"email"`
	Password string `db:"password"`
	Phone string `db:"phone"`
	Origin string `db:"origin"`
	Active bool `db:"active"`
	Image string `db:"image"`
	Role string `db:"role"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}