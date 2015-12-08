package utils

import (
	"agentdesks"
	"agentdesks/models"
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"strconv"
)

/* Checks and logs the error
   Writes back the Error as HTTP response
*/

type ErrHandle func(http.ResponseWriter, *http.Request, httprouter.Params) error

func (eh ErrHandle) ToHandle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := eh(w, r, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Checks and logs the error
func PanicError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s : %s", msg, err))
	}
}

func InitializeErrorLog(config *agentdesks.Config) *os.File {
	logFile, err := os.OpenFile(config.Logs.ErrorLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	PanicError(err, "Error opening access log")
	log.SetOutput(logFile)
	log.Println("Server restarted !! Begining of error log !!")
	return logFile
}

func GetErrorContent(code int) (resp *models.HTTPErrorResponse) {
	var buffer bytes.Buffer
	buffer.WriteString("RUBY_")
	buffer.WriteString(strconv.Itoa(code))
	return &models.HTTPErrorResponse{
		ErrorCode: buffer.String(),
		Message:   ErrorCodes[code],
	}
}
