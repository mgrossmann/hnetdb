package users

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type UserLogin struct {
	User User `json:"user"`
}

type LoggedInUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type UserLoginHandler struct {
	Path           string
	UserRepository UserRepository
}

func (u *UserLoginHandler) Login(writer http.ResponseWriter, request *http.Request) {
	requestBody, _ := ioutil.ReadAll(request.Body)
	userLoginRequest := UserLogin{}
	_ = json.Unmarshal(requestBody, &userLoginRequest)
	requestUser := userLoginRequest.User
	user, _ := u.UserRepository.FindByEmailAndPassword(
		requestUser.Email,
		requestUser.Password)

	if user == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, _ := CreateToken(user)

	writer.WriteHeader(http.StatusOK)
	responseBody := UserLogin{
		User: User{
			Username: user.Username,
			Email:    user.Email,
			Token:    token,
		}}
	bytes, _ := json.Marshal(&responseBody)
	_, _ = writer.Write(bytes)
}

func CreateToken(user *User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_ACCESS")))
}
