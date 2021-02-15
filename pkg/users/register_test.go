package users_test

import (
	"github.com/mvslovers/hnetdb/pkg/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type FakeUserRepository struct {

}

func (FakeUserRepository) RegisterUser(user users.User) error {
	return nil
}

var _ = Describe("User", func() {

	It("should register", func() {

		handler := users.UserHandler{
			Path:           "/users",
			UserRepository: &FakeUserRepository{},
		}

		testResponseWriter := httptest.NewRecorder()
		requestBody := strings.NewReader("{\"user\":{\"email\": \"mig@dearn\",\"username\": \"mig\",\"password\": \"s3cr3t\"}}")

		handler.Register(testResponseWriter, httptest.NewRequest(http.MethodPost,"/users", requestBody))

		Expect(testResponseWriter.Code).To(Equal(201))

		responseBody, _ := ioutil.ReadAll(testResponseWriter.Body)
		Expect(string(responseBody)).To(Equal("{\"user\":{\"username\":\"mig\",\"email\":\"mig@dearn\"}}"))

	})

})