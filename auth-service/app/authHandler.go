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
