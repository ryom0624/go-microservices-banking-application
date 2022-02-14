package domain

import (
	"encoding/json"
	"local.packages/lib/logger"
	"net/http"
	"net/url"
)

type AuthRepository interface {
	IsAuthorized(token string, routeName string, vars map[string]string) bool
}

type RemoteAuthRepository struct{}

func (r RemoteAuthRepository) IsAuthorized(token string, routeName string, vars map[string]string) bool {
	u := buildVerifyURL(token, routeName, vars)

	response, err := http.Get(u)
	if err != nil {
		logger.Error("Error whiling sending request auth server")
		return false
	}
	m := map[string]bool{}
	if err := json.NewDecoder(response.Body).Decode(&m); err != nil {
		logger.Error("Error whiling decoding response from auth server")
		return false
	}

	return m["isAuthorized"]
}

func buildVerifyURL(token string, routeName string, vars map[string]string) string {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost:8181",
		Path:   "/auth/verify",
	}
	q := u.Query()
	q.Add("token", token)
	q.Add("routeName", routeName)
	for k, v := range vars {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func NewAuthRepository() RemoteAuthRepository {
	return RemoteAuthRepository{}
}
