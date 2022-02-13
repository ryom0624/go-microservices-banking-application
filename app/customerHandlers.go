package app

import (
	"github.com/gorilla/mux"
	"go-microservices-banking-application/service"
	"net/http"
)

type Customer struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
}

type CustomerHandler struct {
	service service.CustomerService
}

func (h *CustomerHandler) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	status := q.Get("status")

	customers, err := h.service.GetAllCustomers(status)
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
		return
	}

	writeResponse(w, http.StatusOK, customers)
}

func (h *CustomerHandler) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["customer_id"]

	customer, err := h.service.GetCustomer(id)
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, customer)
	}
}
