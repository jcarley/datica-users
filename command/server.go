package command

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jcarley/datica-users/api/authentication"
	"github.com/jcarley/datica-users/api/public"
	"github.com/jcarley/datica-users/api/users"
	"github.com/jcarley/datica-users/models"
	"github.com/jcarley/datica-users/models/postgres"
	"github.com/jcarley/datica-users/web"
	"github.com/mitchellh/cli"
)

type ServerCommand struct {
	Ui cli.Ui
}

func (this *ServerCommand) Run(args []string) int {

	var port string

	flags := FlagSet("server")
	flags.StringVar(&port, "port", "3000", "port on which HTTP will listen")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Output the header that the server has started
	this.Ui.Output("==> Datica-users server started! Log data will stream in below:\n")

	// how to get the root
	file, err := os.Open("db.hcl")
	if err != nil {
		log.Fatal("Failed to read db.hcl")
		return 1
	}
	databaseInfo, err := postgres.ReadConfig(file)
	if err != nil {
		log.Fatal("Failed to parse db.hcl: ", err.Error())
		return 1
	}
	postgres.Connect(databaseInfo)
	defer postgres.Disconnect()

	server := web.NewHttpServer()
	server.Port = port
	this.RegisterControllers(server)

	log.Printf("Listening on port %s\n", port)

	if err := server.Listen(); err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}

	log.Println("Closing down")

	return 0
}

func (this *ServerCommand) Help() string {
	helpText := `
Usage: datica-users server

	Starts up the datica-users server process.
`
	return strings.TrimSpace(helpText)
}

func (this *ServerCommand) Synopsis() string {
	return "Starts up the datica-users server process."
}

func (this *ServerCommand) RegisterControllers(server *web.HttpServer) {

	userProvider := postgres.NewPostgresUserProvider()
	userRepository := models.NewUserRepository(userProvider)

	healthCheckController := public.NewHealthCheckController()
	healthCheckController.Register(server)

	authenticateController := authentication.NewAuthenticationController(userRepository)
	authenticateController.Register(server)

	usersController := users.NewUsersController(userRepository)
	usersController.Register(server)

	// this is a keep-alive check
	server.Get("/ping", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "pong")
	}))
}
