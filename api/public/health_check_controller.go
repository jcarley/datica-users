package public

import (
	"net/http"

	"github.com/jcarley/datica-users/helper/jsonutil"
	"github.com/jcarley/datica-users/web"
)

type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (this *HealthCheckController) Register(server *web.HttpServer) {
	server.Get("/", this.Show)
}

func (this *HealthCheckController) Show(rw http.ResponseWriter, req *http.Request) {
	greeting := jsonutil.NewValues()
	greeting.Add("greeting", "Hello World")
	greeting.EncodeToWriter(rw)
}
