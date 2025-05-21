package httphandlers

import (
    "net/http"

    "github.com/gorilla/mux"
)

func SetupRouter(userHandler *UserHandler) *mux.Router {
    r := mux.NewRouter()

    authRouter := r.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/signup", userHandler.SignUpHandler).Methods("POST")


    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello from my Go project!"))
    })

    return r
}