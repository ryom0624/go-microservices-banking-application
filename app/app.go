package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go-microservices-banking-application/domain"
	"go-microservices-banking-application/service"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Start() {

	sanityCheck()

	// mux := http.NewServeMux()
	router := mux.NewRouter()

	dbClient := getDBClient()

	customerRepositoryDB := domain.NewCustomerRepositoryDB(dbClient)
	accountRepositoryDB := domain.NewAccountRepositoryDB(dbClient)

	// ch := CustomerHandler{service: service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	ch := CustomerHandler{service: service.NewCustomerService(customerRepositoryDB)}
	ah := AccountHandler{service: service.NewAccountService(accountRepositoryDB)}

	router.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customer/{customer_id:[0-9]+}", ch.getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customer/{customer_id:[0-9]+}/account", ah.NewAccount).Methods(http.MethodPost)
	router.HandleFunc("/customer/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ah.MakeTransaction).Methods(http.MethodPost)

	// exercise
	// router.HandleFunc("/api/time", getCurrentTime).Methods(http.MethodGet)

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func sanityCheck() {
	if os.Getenv("SERVER_ADDRESS") == "" || os.Getenv("SERVER_PORT") == "" {
		log.Fatal("Environment variable not defined...")
	}

}

func getDBClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWORD")
	dbAddr := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
	}
}

func getCurrentTime(w http.ResponseWriter, r *http.Request) {
	tz := r.URL.Query().Get("tz")

	separatedTz := strings.Split(tz, ",")

	if len(separatedTz) == 0 {
		now, _ := time.ParseInLocation(time.RFC3339, time.Now().In(time.UTC).Format(time.RFC3339), time.UTC)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			CurrentTime string `json:"current_time"`
		}{
			CurrentTime: now.String(),
		})

		return

	} else {

		currentTimeByTimezone := make(map[string]string, len(separatedTz))
		for _, t := range separatedTz {
			loc, err := time.LoadLocation(t)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "invalid timezone")
				return
			}
			now, _ := time.ParseInLocation(time.RFC3339, time.Now().In(loc).Format(time.RFC3339), loc)
			currentTimeByTimezone[t] = now.String()
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentTimeByTimezone)

		return
	}

}
