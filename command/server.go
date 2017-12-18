package command

import (
	"log"
	"strings"

	"github.com/jcarley/datica-users/api/public"
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

	// file, err := os.Open("db.hcl")
	// if err != nil {
	// log.Fatal("Failed to read db.hcl")
	// return 1
	// }
	// databaseInfo, err := postgres.ReadConfig(file)
	// if err != nil {
	// log.Fatal("Failed to parse db.hcl: ", err.Error())
	// return 1
	// }
	// postgres.Connect(databaseInfo)
	// defer postgres.Disconnect()

	// this should probably pull items from Tapprd Config
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

	healthCheckController := public.NewHealthCheckController()
	healthCheckController.Register(server)

	// userProvider := postgres.NewPostgresUserProvider()
	// userRepository := models.NewUserRepository(userProvider)

	// signatureController := signatures.NewSignatureController()
	// signatureController.Register(server)

	// authenticateController := authentication.NewAuthenticateController()
	// authenticateController.SignIn = userRepository.Verify
	// authenticateController.LookupUser = userRepository.FindByUsername
	// authenticateController.Register(server)

	// usersController := users.NewUsersController()
	// usersController.CreateUser = userRepository.CreateUser
	// usersController.FindUser = userRepository.FindUser
	// usersController.Register(server)

	// updateController := updates.NewUpdateController()
	// updateController.Register(server)

	// // this is a stand-in for a health check.  Will move this to a proper health check controller
	// server.Get("/ping", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	// fmt.Fprintf(rw, "pong")
	// }))
}
