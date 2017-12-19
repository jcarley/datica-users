package users

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jcarley/datica-users/helper/jsonutil"
	"github.com/jcarley/datica-users/models"
	"github.com/jcarley/datica-users/test"
	"github.com/jcarley/datica-users/web"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserSignup_Integration(t *testing.T) {

	Convey("A user can signup", t, func() {
		server := createUserController()

		userData := jsonutil.NewValues()
		userData.Add("email", "jeff.carley@example.com")
		userData.Add("name", "Jeff Carley")
		userData.Add("password", "secret")
		json, _ := userData.Encode()
		body := bytes.NewBuffer(json)

		req, _ := http.NewRequest("POST", "/user", body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(json)))

		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		So(rw.Code, ShouldEqual, http.StatusCreated)
		So(rawData["user_id"], ShouldEqual, "1234567890")
		So(rawData["email"], ShouldEqual, "jeff.carley@example.com")
		So(rawData["name"], ShouldEqual, "Jeff Carley")
		So(rawData["password"], ShouldBeNil)
		So(rawData["salt"], ShouldBeNil)
	})

}

func TestGetUserData_Integration(t *testing.T) {
	Convey("Get a given users data", t, func() {
		server := createUserController()

		userId, username := "1234567890", "jeff.carley@example.com"
		token := authToken(userId, username)

		url := fmt.Sprintf("/user/%s", username)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("bearer %v", token))

		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		So(rw.Code, ShouldEqual, http.StatusOK)
		So(rawData, ShouldNotBeNil)
	})

	Convey("Returns and error when token and username don't match", t, func() {
		server := createUserController()

		userId, username := "1234567890", "jeff.carley@example.com"
		token := authToken(userId, username)

		url := fmt.Sprintf("/user/%s", "jcarley")

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("bearer %v", token))

		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		So(rw.Code, ShouldEqual, http.StatusUnauthorized)
		So(rawData, ShouldNotBeNil)
	})
}

func TestUpdatingUserData_Integration(t *testing.T) {
	Convey("Update a given authenticated user's data", t, func() {
		server := createUserController()

		userId, username := "1234567890", "jeff.carley@example.com"
		token := authToken(userId, username)

		url := fmt.Sprintf("/user/%s", username)

		userData := jsonutil.NewValues()
		userData.Add("name", "John Doe")
		json, _ := userData.Encode()
		body := bytes.NewBuffer(json)

		req, _ := http.NewRequest("PUT", url, body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("bearer %v", token))

		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		// even though the UpdateUser api call only allows updating the user's name
		// we assert that email and user_id didn't change.  ensures the correct user
		// object is being passed to the update inside of the controller
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(rawData, ShouldNotBeNil)
		So(rawData["name"], ShouldEqual, "John Doe")
		So(rawData["email"], ShouldEqual, "jeff.carley@example.com")
		So(rawData["user_id"], ShouldEqual, "1234567890")
	})
}

func TestDeleteUserData_Integration(t *testing.T) {
	Convey("Delete a given authenticated user's data", t, func() {
		server := createUserController()

		userId, username := "1234567890", "jeff.carley@example.com"
		token := authToken(userId, username)

		url := fmt.Sprintf("/user/%s", username)

		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("bearer %v", token))

		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		So(rw.Code, ShouldEqual, http.StatusOK)
	})
}

func TestAuthenticatedRequestReturnsUsersInfo_Integration(t *testing.T) {

	Convey("Request for root with a token returns current user info", t, func() {

		server := createUserController()

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

func createUserController() *web.HttpServer {
	server := web.NewHttpServer()

	userRepository := createUserRepository()
	controller := NewUsersController(userRepository)
	controller.Register(server)

	return server
}

func createUserRepository() *models.UserRepository {
	return models.NewUserRepository(&InMemoryUserProvider{})
}

type InMemoryUserProvider struct {
}

func (this *InMemoryUserProvider) AddUser(user models.User) (models.User, error) {
	return models.User{
		UserId: "1234567890",
		Name:   "Jeff Carley",
		Email:  "jeff.carley@example.com",
	}, nil
}

func (this *InMemoryUserProvider) UpdateUser(userId string, user models.User) (models.User, error) {
	return models.User{
		UserId:    userId,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (this *InMemoryUserProvider) DeleteUser(userId string) error {
	return nil
}

func (this *InMemoryUserProvider) FindByUsername(username string) (models.User, error) {
	return models.User{UserId: "1234567890", Name: "Jeff Carley", Email: username}, nil
}

func (this *InMemoryUserProvider) FindUser(userId string) (models.User, error) {
	return models.User{}, nil
}
