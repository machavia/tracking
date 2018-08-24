package validator

import (
	"../models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/schema"
	"reflect"
	"strconv"
	"time"
)

//client id must be a 6 integer characters
func ClientId(id string) bool {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		println(err)
		return false
	}
	if len(id) == 5 && idInt >= 10000 && idInt <= 99999 {
		return true
	}

	return false
}

//visitor id must be 17 characters long integer
func VisitorId(id string) bool {

	if len(id) == 17 {
		return true
	}
	return false
}

//validation of the http request received
//we need to make sure the data that enter the data-pipeline are clean and safe
//Then we jsonify it
func validateStructAndJsonify(hitStruct interface{}, form map[string][]string) ([]byte, error) {

	var empty []byte
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(hitStruct, form); err != nil {
		fmt.Println(err)
		return empty, errors.New("Decoder error")
	}

	var currentTime uint64 = uint64(time.Now().Unix())

	reflect.ValueOf(hitStruct).Elem().FieldByName("Time").SetUint(currentTime)

	hitJson, err := json.Marshal(hitStruct)
	if err != nil {
		fmt.Println(err)
		return empty, errors.New("Json error")
	}

	return hitJson, nil

}

//each http request (hit) has a declared type which we need to know for validation purpose
func GetHitType(HitType string, form map[string][]string) ([]byte, error) {

	var myjson []byte
	var err error

	switch HitType {
	case "pag":
		var hitmodel = new(models.Hit24Pag)
		myjson, err = validateStructAndJsonify(hitmodel, form)

	default:
		return myjson, nil
	}

	if err != nil {
		return myjson, err
	}
	return myjson, nil

}
