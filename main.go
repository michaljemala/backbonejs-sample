package main

import (
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type (
	Session struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Date  string `json:"date"`
	}
)

var (
	storage = make(map[int]Session)
)

func init() {
	// Sample data
	storage[1] = Session{1, "Sess1", "06/02/2013"}
	storage[2] = Session{2, "Sess2", "07/02/2013"}
	storage[3] = Session{3, "Sess3", "08/02/2013"}
	storage[4] = Session{4, "Sess4", "11/02/2013"}
	storage[5] = Session{5, "Sess5", "12/02/2013"}
}

func main() {
	r := mux.NewRouter()

	// Handle /schedule
	r.HandleFunc("/schedule", ListHandler).Methods("GET")
	r.HandleFunc("/schedule", CreateHandler).Methods("POST")
	r.HandleFunc("/schedule/{id:[0-9]+}", GetHandler).Methods("GET")
	r.HandleFunc("/schedule/{id:[0-9]+}", UpdateHandler).Methods("PUT")
	r.HandleFunc("/schedule/{id:[0-9]+}", DeleteHandler).Methods("DELETE")

	// Handle resources
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("www")))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ListHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")

	s := []Session{}
	for _, v := range storage {
		s = append(s, v)
	}

	enc := json.NewEncoder(rw)
	enc.Encode(&s)
}

func CreateHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
}

func GetHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
}

func UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
}

func DeleteHandler(rw http.ResponseWriter, r *http.Request) {
	id_str := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(id_str, 0, 0)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	delete(storage, int(id))

	rw.WriteHeader(http.StatusOK)
}
