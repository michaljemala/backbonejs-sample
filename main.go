package main

import (
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

type (
	Session struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Date  string `json:"date"`
	}
	Storage struct {
		sync.RWMutex
		data map[int]Session
	}
)

var (
	storage  = Storage{data: make(map[int]Session)}
	sequence = sequenceGenerator(1)
)

func init() {
	// Pre-generate sample data set
	var idx int
	for i := 1; i <= 5; i++ {
		idx = <-sequence
		storage.data[idx] = Session{idx, fmt.Sprintf("Sess%d", idx), "06/02/2013"}
	}
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
	log.Fatal(http.ListenAndServe(":8080", logger(http.DefaultServeMux)))
}

func ListHandler(rw http.ResponseWriter, r *http.Request) {
	storage.RLock()
	var keys []int
	for k := range storage.data {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	sessions := []Session{}
	for _, k := range keys {
		sessions = append(sessions, storage.data[k])
	}
	storage.RUnlock()

	rw.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(rw)
	enc.Encode(&sessions)
}

// Handlers
func CreateHandler(rw http.ResponseWriter, r *http.Request) {
	var s Session

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		data, _ := json.Marshal(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(data)
		return
	}

	s.Id = <-sequence

	storage.Lock()
	storage.data[s.Id] = s
	storage.Unlock()

	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(rw)
	enc.Encode(&s)
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

	storage.Lock()
	delete(storage.data, int(id))
	storage.Unlock()

	rw.WriteHeader(http.StatusNoContent)
}

// Utils
func sequenceGenerator(start int) <-chan int {
	c := make(chan int)
	go func(s int) {
		for i := s; ; i++ {
			c <- i
		}
	}(start)
	return c
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(wr, r)
	})
}
