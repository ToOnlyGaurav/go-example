package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gorilla/mux"
)

type Rest struct {
	router *mux.Router
}

func (r *Rest) start() {
	fmt.Println("Started..")
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:               10000,
		MaxConcurrentRequests: 1000,
		ErrorPercentThreshold: 25,
	})
	r.router = mux.NewRouter()
	r.initializeRoutes()
}

func (c *Rest) initializeRoutes() {
	c.router.HandleFunc("/v1/config", c.handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":1234", c.router))
}
func (c *Rest) handler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, c.getMessage())
}

func (c *Rest) getRandomSleep() int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1000) // n will be between 0 and 10
	fmt.Printf("Sleeping %d millisecond...\n", n)
	return n
}
func (c *Rest) getMessage() interface{} {
	hystrix.Go("my_command", func() error {
		// talk to other services
		time.Sleep(time.Duration(c.getRandomSleep() * int(time.Millisecond)))
		return nil
	}, nil)

	return "Hello " + time.Now().String()
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
