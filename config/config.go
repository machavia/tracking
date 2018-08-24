package config

import (
	"encoding/json"
	"os"
)

var (
	Config *Configuration
)

type Configuration struct {
	AwsKey             string `json:"aws_key"`
	AwsSecret          string `json:"aws_secret"`
	RedisHost          string `json:"redis_host"`
	KinesisStreamName  string `json:"kinesis_stream"`
	FirehoseStreamName string `json:"firehose_stream"`
	ApiHost            string `json:"api_host"`
	ApiKey             string `json:"api_key"`
}

func Load() {

	//filename is the path to the json config file
	configuration := Configuration{}
	pwd, _ := os.Getwd() //to find the correct pass in AWS EBS
	file, err := os.Open(pwd + "/config/config.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	Config = &configuration
}
