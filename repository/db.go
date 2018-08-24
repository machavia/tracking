package repository

import "strconv"

//Get a visitor buckets (as a string) from Redis
func GetVisitor(clientId uint32, visitorId uint64) string {

	//the redis keys are built with the client id and visitor id to make sure they are unique
	redisKey := strconv.Itoa(int(clientId)) + ":" + strconv.FormatUint(visitorId, 10)

	val, err := utils.RedisClient.HGetAll(redisKey).Result()
	if err != nil {
		panic(err)
	}

	//if the visitor has the bucket 998 he is intent
	if val["bucket_partner_1"] == "998" {
		//a little crappy here, but the redis hash keys will complety change shortly so no need to spend time here
		return "{\"1\":\"998\"}"
	}

	//if the visitor has no buckets
	return "{}"

}
