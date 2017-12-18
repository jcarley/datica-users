package authentication

import (
	"errors"
	"net/http"

	"github.com/jcarley/datica-users/helper/jsonutil"
	"github.com/jcarley/datica-users/models"
	"github.com/jcarley/datica-users/web"
)

var (
	ErrAuthenticationFailed = errors.New("Authentication for user failed")
)

type CredentialProvider interface {
	SignIn(username string, password string) (*models.User, error)
}

type AuthenticationController struct {
	credentialProvider CredentialProvider
}

func NewAuthenticationController(credentialProvider CredentialProvider) *AuthenticationController {
	return &AuthenticationController{
		credentialProvider,
	}
}

func (this *AuthenticationController) Register(server *web.HttpServer) {
	server.Post("/auth", this.Authenticate)
}

func (this *AuthenticationController) Authenticate(rw http.ResponseWriter, req *http.Request) {

	var json map[string]interface{}
	jsonutil.DecodeJSONFromReader(req.Body, &json)
	req.Body.Close()

	username := json["username"].(string)
	password := json["password"].(string)

	if user, err := this.credentialProvider.SignIn(username, password); err == nil {

		claim := web.NewClaim(user.UserId, user.Email)
		sig, token := web.Signature(claim)
		authToken := web.NewAuthToken(token, sig)

		jsonutil.EncodeJSONToWriter(rw, &authToken)

	} else {
		web.Error(rw, ErrAuthenticationFailed, http.StatusUnauthorized)
	}
}
