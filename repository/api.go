package repository

import (
	"../utils"
	"time"
)

const ClientListRetentionTime = 60 * 60

//getting client status from the hot memory
//If the statuses have not been updated for a while, we're updating them from the Stormize API
func GetClientsStatus(clientId uint32) (bool, bool) {

	//if the client list is no longer updated
	if utils.LastClientsUpdate == 0 || time.Now().Unix()-utils.LastClientsUpdate > ClientListRetentionTime {
		//getting the list in the API
		//Function called in async mode to avoid blocking (http) client serving
		go utils.GetClientsList()
	}

	//if the client exist we return the statuses
	if _, ok := utils.Clients[clientId]; ok {
		return utils.Clients[clientId]["tracking"], utils.Clients[clientId]["prediction"]
	}

	//else all statuses to disabled
	return false, false
}
