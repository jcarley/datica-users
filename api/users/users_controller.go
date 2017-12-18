package users

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jcarley/datica-users/helper/jsonutil"
	"github.com/jcarley/datica-users/models"
	"github.com/jcarley/datica-users/web"
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
	server.SecureGet("/", this.Show)
}

func (this *UsersController) Show(rw http.ResponseWriter, req *http.Request) {
	token := req.Context().Value("user")
	claim := token.(*jwt.Token).Claims.(jwt.MapClaims)
	username := claim["username"].(string)

	if user, err := this.userRepository.FindByUsername(username); err == nil {
		user.Salt = ""
		user.Password = ""
		jsonutil.EncodeJSONToWriter(rw, &user)
	}
}
