package utils

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func WrapResponse(writer http.ResponseWriter, content interface{}, status int) {
	writer.Header().Set(FORMAT, "application/json")
	writer.WriteHeader(status)
	if content != nil {
		responseJson, err := json.Marshal(content)
		PanicError(err, "Error wrapping response")
		writer.Write(responseJson)
	}
}

//Util to convert HTTPRouter handler to net/http handler
func WrapHandler(hr func(http.ResponseWriter, *http.Request)) httprouter.Handle {
	h := alice.New(context.ClearHandler).ThenFunc(hr)
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}
