package app

import (
	"auth/dto"
	"auth/service"
	"encoding/json"
	"local.packages/lib/logger"
	"net/http"
)

type authHandler struct {
	service service.AuthService
}

func (h authHandler) Login(w http.ResponseWriter, req *http.Request) {

	var request dto.LoginRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		logger.Error("Error while decoding login request " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, appError := h.service.Login(request)
	if appError != nil {
		writeResponse(w, appError.Code, appError.AsMessage())
	}
	writeResponse(w, http.StatusOK, token)
}

func (h authHandler) NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusNotImplemented, "not implement")
}

func (h authHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] == "" {
		writeResponse(w, http.StatusForbidden, notAuthorizedResponse("missing token"))
		return
	}
	appErr := h.service.Verify(urlParams)
	if appErr != nil {
		writeResponse(w, appErr.Code, notAuthorizedResponse(appErr.Message))
		return
	}
	writeResponse(w, http.StatusOK, authorizedResponse())

}

func authorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}
}

func notAuthorizedResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      message,
	}
}
