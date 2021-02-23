package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UserRegistration struct {
	User User `json:"user"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token"`
}

type UserRegistrationHandler struct {
	Path           string
	UserRepository UserRepository
}

func (u *UserRegistrationHandler) Register(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestBody, _ := ioutil.ReadAll(request.Body)
	userRegistrationRequest := UserRegistration{}
	err := json.Unmarshal(requestBody, &userRegistrationRequest)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	requestUser := userRegistrationRequest.User
	err = u.UserRepository.RegisterUser(&requestUser)
	if err != nil {

		panic(err)
	}

	writer.WriteHeader(201)
	writer.Header().Add("Content-Type", "application/json")
	userRegistrationResponse := UserRegistration{
		User: User{
			Username: requestUser.Username,
			Email:    requestUser.Email,
		}}
	bytes, _ := json.Marshal(&userRegistrationResponse)
	_, _ = writer.Write(bytes)
}
