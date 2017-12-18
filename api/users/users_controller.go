package users

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jcarley/datica-users/helper/jsonutil"
	"github.com/jcarley/datica-users/models"
	"github.com/jcarley/datica-users/web"
)

var (
	ErrInvalidAuthToken = errors.New("Auth token is missing or invalid")
)

type UsersController struct {
	userRepository *models.UserRepository
}

func NewUsersController(userRepository *models.UserRepository) *UsersController {
	return &UsersController{
		userRepository,
	}
}

func (this *UsersController) Register(server *web.HttpServer) {
	server.Post("/user", this.Create)
	server.SecureGet("/", this.Show)
}

func (this *UsersController) Create(rw http.ResponseWriter, req *http.Request) {
}

func (this *UsersController) Show(rw http.ResponseWriter, req *http.Request) {
	if token, ok := web.TokenFromContext(req.Context()); ok {
		claim := token.Claims.(jwt.MapClaims)
		username := claim["username"].(string)

		if user, err := this.userRepository.FindByUsername(username); err == nil {
			user.Salt = ""
			user.Password = ""
			jsonutil.EncodeJSONToWriter(rw, &user)
		}
	} else {
		web.Error(rw, ErrInvalidAuthToken, http.StatusUnauthorized)
	}
}
