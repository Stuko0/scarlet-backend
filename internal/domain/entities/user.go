package entities

import (
	"encoding/json"
	"net/http"
)

type User  struct{
	Id int	`json:"id"`
	Name string `json:"name"`
	Lastname string `json:"lastname"`
	Email string `json:"email"`
	Pwd string `json:"pwd"`
}

var (
	users []User
)

func init(){
	users =  []User{{Id: 1, Name: "Alex", Lastname: "Villanueva", Email: "avplaying@gmail.com", Pwd: "1234a"}}
}

func GetUsers(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	result, err := json.Marshal(users)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error marshalling the users array"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write(result)
}

func AddUsers(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var user User
	err:=json.NewDecoder(req.Body).Decode(&user)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error marshalling the request"}`))
		return
	}
	user.Id=len(users)+1
	users=append(users, user)
	resp.WriteHeader(http.StatusOK)
	result, err := json.Marshal(users)
	resp.Write(result)
}