package authentication

import (
	"bytes"
	"errors"
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

func TestAuthenticateReturnsAuthToken_Integration(t *testing.T) {

	server := web.NewHttpServer()

	signInFunc := CredentialProviderHandlerFunc(func(username string, password string) (*models.User, error) {
		return &models.User{
			UserId: "1234567890",
			Name:   "Jeff Carley",
			Email:  username,
		}, nil
	})

	controller := NewAuthenticationController(signInFunc)
	controller.Register(server)

	Convey("Auth endpoint returns an auth token", t, func() {

		data := jsonutil.NewValues()
		data.Set("username", "jeff.carley@example.com")
		data.Set("password", "secret")
		json, _ := data.Encode()
		body := bytes.NewBuffer(json)

		req, _ := http.NewRequest("POST", "/auth", body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(json)))

		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		authToken := web.AuthToken{
			Header: make(map[string]interface{}),
		}
		test.DecodeStruct(t, &authToken, rawData)

		So(rw.Code, ShouldEqual, http.StatusOK)
		So(authToken.Header, ShouldNotBeNil)
		So(authToken.Claim, ShouldNotBeNil)
		So(authToken.Claim.UserId, ShouldEqual, "1234567890")
		So(authToken.Claim.Username, ShouldEqual, "jeff.carley@example.com")
		So(authToken.Signature, ShouldNotBeNil)
	})
}

func TestAuthenticateFailure_Integration(t *testing.T) {

	server := web.NewHttpServer()

	signInFunc := CredentialProviderHandlerFunc(func(username string, password string) (*models.User, error) {
		return &models.User{}, errors.New("Unknow user")
	})

	controller := NewAuthenticationController(signInFunc)
	controller.Register(server)

	Convey("Auth endpoint returns a 401 on auth failure", t, func() {

		data := jsonutil.NewValues()
		data.Set("username", "jeff.carley@example.com")
		data.Set("password", "wrongpassword")
		json, _ := data.Encode()
		body := bytes.NewBuffer(json)

		req, _ := http.NewRequest("POST", "/auth", body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(json)))

		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		var errorMsg web.ErrorMessage
		test.DecodeStruct(t, &errorMsg, rawData["error"])

		So(rw.Code, ShouldEqual, http.StatusUnauthorized)
		So(errorMsg.Status, ShouldEqual, http.StatusUnauthorized)
	})
}

type CredentialProviderHandlerFunc func(username string, password string) (*models.User, error)

func (h CredentialProviderHandlerFunc) SignIn(username string, password string) (*models.User, error) {
	return h(username, password)
}

// {
// "header": {
// "typ": "JWT",
// "alg": "HS256"
// },
// "claim": {
// "userid": 1
// },
// "signature": ""
// }
