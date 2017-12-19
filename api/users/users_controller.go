package users

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jcarley/datica-users/helper/encoding"
	"github.com/jcarley/datica-users/helper/jsonutil"
	"github.com/jcarley/datica-users/helper/structutil"
	"github.com/jcarley/datica-users/models"
	"github.com/jcarley/datica-users/web"
)

var (
	ErrCreateUser       = errors.New("Failed to create user")
	ErrDeleteUser       = errors.New("Deleting user failed")
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
	server.SecureGet("/", this.Index)
	server.SecureGet("/user/{username}", this.Show)
	server.SecurePut("/user/{username}", this.Update)
	server.SecureDelete("/user/{username}", this.Delete)
}

func (this *UsersController) Index(rw http.ResponseWriter, req *http.Request) {
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

func (this *UsersController) Show(rw http.ResponseWriter, req *http.Request) {

	if token, ok := web.TokenFromContext(req.Context()); ok {
		vars := mux.Vars(req)
		routeUsername := vars["username"]

		claim := token.Claims.(jwt.MapClaims)
		tokenUsername := claim["username"].(string)

		if routeUsername == tokenUsername {
			if user, err := this.userRepository.FindByUsername(tokenUsername); err == nil {
				user.Salt = ""
				user.Password = ""
				jsonutil.EncodeJSONToWriter(rw, &user)
			}
		} else {
			web.Error(rw, ErrInvalidAuthToken, http.StatusUnauthorized)
		}
	}

}

func (this *UsersController) Create(rw http.ResponseWriter, req *http.Request) {

	var json map[string]interface{}
	jsonutil.DecodeJSONFromReader(req.Body, &json)
	req.Body.Close()

	email := json["email"].(string)
	name := json["name"].(string)
	password := json["password"].(string)

	salt := encoding.GenerateSalt()
	encryptedPassword := encoding.Key(password, salt)

	user := models.User{
		Email:    email,
		Name:     name,
		Salt:     salt,
		Password: encryptedPassword,
	}

	if newUser, err := this.userRepository.AddUser(user); err == nil {
		newUser.Password = ""
		newUser.Salt = ""
		rw.WriteHeader(http.StatusCreated)
		jsonutil.EncodeJSONToWriter(rw, &newUser)
	} else {
		web.Error(rw, ErrCreateUser, http.StatusInternalServerError)
	}

}

func (this *UsersController) Update(rw http.ResponseWriter, req *http.Request) {
	if token, ok := web.TokenFromContext(req.Context()); ok {
		vars := mux.Vars(req)
		routeUsername := vars["username"]

		claim := token.Claims.(jwt.MapClaims)
		tokenUsername := claim["username"].(string)
		tokenUserId := claim["userid"].(string)

		if routeUsername == tokenUsername {

			var json map[string]interface{}
			jsonutil.DecodeJSONFromReader(req.Body, &json)
			req.Body.Close()

			var updateUser models.User
			structutil.Decode(&updateUser, json)

			// Pulling the user record first isn't great.  sticking with this
			// implementation for interest in time.
			currentUser, err := this.userRepository.FindByUsername(tokenUsername)
			if err != nil {
				// didn't find user.  what status code to return for this?
				web.Error(rw, err, http.StatusUnprocessableEntity)
				return
			}

			// current implementation only allows changing the name of the user.  changing the email
			// which is also acts as the username for signin could result in a loss of login or a unique
			// constraint database error.  Further implementation could allow an admin to perform this
			currentUser.Name = updateUser.Name

			// we got the currentUser object from the database and updated its attributes.  we need to
			// make sure we send the currentUser object to the UpdateUser function
			if user, err := this.userRepository.UpdateUser(tokenUserId, currentUser); err == nil {
				user.Salt = ""
				user.Password = ""
				jsonutil.EncodeJSONToWriter(rw, &user)
			} else {
				web.Error(rw, err, http.StatusInternalServerError)
			}
		} else {
			web.Error(rw, ErrInvalidAuthToken, http.StatusUnauthorized)
		}

	}
}

func (this *UsersController) Delete(rw http.ResponseWriter, req *http.Request) {
	if token, ok := web.TokenFromContext(req.Context()); ok {
		vars := mux.Vars(req)
		routeUsername := vars["username"]

		claim := token.Claims.(jwt.MapClaims)
		tokenUsername := claim["username"].(string)
		tokenUserId := claim["userid"].(string)

		if routeUsername == tokenUsername {
			if err := this.userRepository.DeleteUser(tokenUserId); err != nil {
				web.Error(rw, ErrDeleteUser, http.StatusInternalServerError)
			}
		} else {
			web.Error(rw, ErrInvalidAuthToken, http.StatusUnauthorized)
		}

	}
}
