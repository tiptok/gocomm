package httpx

import (
	"encoding/json"
	"fmt"
	"github.com/tiptok/gocomm/pkg/log"
	"net/http"
)

func Error(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func OkJson(w http.ResponseWriter, v interface{}) {
	WriteJson(w, http.StatusOK, v)
}

func WriteJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)

	if bs, err := json.Marshal(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			log.Error("write response failed, error: ", err.Error())
		}
	} else if n < len(bs) {
		log.Error(fmt.Sprintf("actual bytes: %d, written bytes: %d", len(bs), n))
	}
}