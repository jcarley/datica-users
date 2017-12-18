package web

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type key int

const (
	RequestVarsKey key = 0
)

func NewContextWithRequestVars(ctx context.Context, req *http.Request) context.Context {

	// Reads the route variables for the current request, if any.
	vars := mux.Vars(req)
	if vars == nil {
		vars = make(map[string]string)
	}

	// Reads the values from the query string.  Values from the body will be overwritten
	values := req.URL.Query()
	for key, _ := range values {
		vars[key] = values.Get(key)
	}

	return context.WithValue(ctx, RequestVarsKey, vars)
}

func VarsFromContext(ctx context.Context) map[string]string {
	return ctx.Value(RequestVarsKey).(map[string]string)
}
