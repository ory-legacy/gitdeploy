package main

import (
	"github.com/ory-am/google-json-style-response/responder"
	"log"
	"net/http"
)

func cleanUpSession(w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Please delete your cookies.")
		return err
	}
	delete(session.Values, sessionCurrentDeployment)
	if err := session.Save(r, w); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return err
	}
	return nil
}

// Set the different CORS headers required for CORS request
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

func responseError(w http.ResponseWriter, code int, message string) {
	log.Printf("Error %d: %s", code, message)
	response := responder.New(ApiVersion)
	response.Write(w, response.Error(code, message))
}

func responseSuccess(w http.ResponseWriter, data interface{}) {
	response := responder.New(ApiVersion)
	response.Write(w, response.Success(data))
}
