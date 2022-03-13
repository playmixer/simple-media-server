package apperror

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(handle appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError

		//Disable CORS
		w = disabledCors(w)

		err := handle(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrPathNotFound) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrPathNotFound.Marshal())
					return
				} else if errors.Is(err, ErrCantMarshal) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrCantMarshal.Marshal())
					return
				}
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write([]byte(err.Error()))
		}
	}
}

func disabledCors(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, X-Auth-Token, Origin")
	return w
}
