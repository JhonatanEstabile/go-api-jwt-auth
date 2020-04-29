package util

import (
	"api-jwt-auth/db"
	"api-jwt-auth/model"
	"strconv"
)

//FetchAuth load the user id cached in redis
func FetchAuth(authD *model.AccessDetails) (uint64, error) {
	client := db.GetRedisClient()
	defer client.Close()

	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseUint(userid, 10, 64)

	return userID, nil
}

//DeleteAuth delete user auth session in redis
func DeleteAuth(givenUUID string) (int64, error) {
	client := db.GetRedisClient()
	defer client.Close()

	deleted, err := client.Del(givenUUID).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
