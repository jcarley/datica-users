package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jcarley/datica-users/models"
	"github.com/jcarley/datica-users/test"
	"github.com/jcarley/datica-users/web"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthenticatedRequestReturnsUsersInfo_Integration(t *testing.T) {

	server := web.NewHttpServer()

	userRepository := createUserRepository()
	controller := NewUsersController(userRepository)
	controller.Register(server)

	Convey("Request for root with a token returns current user info", t, func() {

		userId, username := "1234567890", "jeff.carley@example.com"
		token := authToken(userId, username)

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("bearer %v", token))

		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		So(rw.Code, ShouldEqual, http.StatusOK)
		So(rawData, ShouldNotBeNil)
		So(rawData["user_id"], ShouldEqual, userId)
		So(rawData["email"], ShouldEqual, username)
	})
}

func authToken(userId, username string) string {
	claim := web.NewClaim(userId, username)
	sig, _ := web.Signature(claim)
	return sig
}

func createUserRepository() *models.UserRepository {
	return models.NewUserRepository(&InMemoryUserProvider{})
}

type InMemoryUserProvider struct {
}

func (this *InMemoryUserProvider) AddUser(user models.User) (models.User, error) {
	return models.User{}, nil
}

func (this *InMemoryUserProvider) FindByUsername(username string) (models.User, error) {
	return models.User{UserId: "1234567890", Name: "Jeff Carley", Email: username}, nil
}

func (this *InMemoryUserProvider) FindUser(userId string) (models.User, error) {
	return models.User{}, nil
}
