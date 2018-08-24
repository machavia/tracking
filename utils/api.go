package utils

import (
	"../config"
	"../models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	LastClientsUpdate int64 = 0
	Clients           map[uint32]map[string]bool
)

//Download from the Stormize Api all the active clients
//Store the prediction and the tracking status per client in (package var) Clients map
func GetClientsList() {
	fmt.Println("Getting visitor list from API")

	//http client with restrictive timeout
	var httpClient = &http.Client{Timeout: 10 * time.Second}

	//the api end point is shared by the Config package along with the API key
	url := config.Config.ApiHost + "client/getstatuslist?key=" + config.Config.ApiKey
	r, err := httpClient.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	//mapping the result to the right structure
	var jsonClients models.ClientsList
	err = json.NewDecoder(r.Body).Decode(&jsonClients)
	if err != nil {
		fmt.Println(err)
	}

	//releasing the body
	r.Body.Close()

	//foreach over the clients to bind to the Clients map with helpfull structure
	var clients = map[uint32]map[string]bool{}
	for _, element := range jsonClients.Payload {

		clients[element.PublicId] = map[string]bool{
			"tracking":   element.Tracking,
			"prediction": element.Prediction}
	}

	//updating the package vars
	Clients = clients

	//LastClientsUpdate helps other packages to know if the list is up to date
	LastClientsUpdate = time.Now().Unix()

	fmt.Println("Clients loaded ready now...")

}
