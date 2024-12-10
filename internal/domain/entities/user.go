package entities

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Psw      string `json:"psw"`
	Origin   string `json:"origin"`
}
