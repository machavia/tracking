package main

import (
	"./config"
	"./controllers"
	"./utils"
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	//loading configuration in /config folder
	config.Load()

	//opening a persistent connection with redis
	utils.OpenDbConnexion()

	//creating authenticated session with aws credentials given by the config package
	utils.CreateAwsSession()

	//getting list of active client we want to authorize
	utils.GetClientsList()

}

func main() {

	//hello word, for health check and load balancing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	//managing the received hits
	http.HandleFunc("/hit", controllers.HandleHit)

	//getting port tu use
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
