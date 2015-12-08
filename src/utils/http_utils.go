package utils

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"log"
	"net/http"
)

func IsValidRequest(r *http.Request, w http.ResponseWriter) bool {
	userAgent := r.Header.Get(USER_AGENT)
	if userAgent == "" || userAgent != "agent-php" {
		log.Println("Invalid request headers")
		WrapResponse(w, GetErrorContent(2), http.StatusBadRequest)
		return false
	}
	return true
}

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
