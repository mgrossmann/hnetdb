package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type RegisterationUser struct {
	User User	`json:"user"`
}

type User struct {
	Username string	`json:"username"`
	Email string	`json:"email"`
	Password string	`json:"password,omitempty"`
}

type UserHandler struct {
	Path string
	UserRepository UserRepository
}

func (u *UserHandler) Register(writer http.ResponseWriter, request *http.Request) {

	registerUserRequest  := RegisterationUser{}
	registerUserResponse := RegisterationUser{}

	requestBody, _ := ioutil.ReadAll(request.Body)

	_ = json.Unmarshal(requestBody, &registerUserRequest)
	requestUser := registerUserRequest.User

	_ = u.UserRepository.RegisterUser(requestUser)

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Add("Content-Type", "application/json")

	responseUser := User{
		Username: requestUser.Username,
		Email:    requestUser.Email,
	}

	registerUserResponse.User = responseUser

	bytes, _ := json.Marshal(&registerUserResponse)
	_, _ = writer.Write(bytes)

}
