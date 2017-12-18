package web

import (
	"net/http"

	"github.com/jcarley/datica-users/helper/jsonutil"
)

type ErrorMessage struct {
	Message    string `json:"message"`
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
}

func Error(w http.ResponseWriter, err error, status int) {
	msg := ErrorMessage{
		Message:    err.Error(),
		Status:     status,
		StatusText: http.StatusText(status),
	}
	errMap := make(map[string]ErrorMessage)
	errMap["error"] = msg
	errJson, _ := jsonutil.EncodeJSONToString(&errMap)
	http.Error(w, errJson, status)
}

// {
// "error": {
// "message": "Some useful message about the error",
// "status": "403",
// "status_text": "Forbidden"
// }
// }
