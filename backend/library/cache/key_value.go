//
// Date: 2/27/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cache

import (
	"encoding/json"
	"go/build"
	"os"
	"time"

	"app.options.cafe/library/services"
	"github.com/go-redis/redis"
	env "github.com/jpfuentes2/go-env"
)

var (
	redisConnection *redis.Client
)

//
// Init - startup Redis
//
func init() {

	// Helpful for testing
	env.ReadEnv(build.Default.GOPATH + "/src/app.options.cafe/.env")

	// Setup Redis connection
	redisConnection = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Make sure we can ping the host.
	_, err := redisConnection.Ping().Result()

	if err != nil {
		services.Fatal(err)
	}
}

//
// Delete a stored key. We do not
// return error as if Redis is not up we should
// have a massive fail.
//
func Delete(key string) {

	// Delete in redis
	err := redisConnection.Del(key).Err()

	if err != nil {
		services.Fatal(err)
	}
}

//
// Store a key value into memory. We do not
// return error as if Redis is not up we should
// have a massive fail.
//
func Set(key string, value interface{}) {

	// Pack up into JSON
	b, err := json.Marshal(&value)

	if err != nil {
		services.Fatal(err)
	}

	// Store in redis
	err = redisConnection.Set(key, b, 0).Err()

	if err != nil {
		services.Fatal(err)
	}
}

//
// Store a key value into memory that expires. We do not
// return error as if Redis is not up we should
// have a massive fail.
//
func SetExpire(key string, expire time.Duration, value interface{}) {

	// Pack up into JSON
	b, err := json.Marshal(&value)

	if err != nil {
		services.Fatal(err)
	}

	// Store in redis
	err = redisConnection.Set(key, b, expire).Err()

	if err != nil {
		services.Fatal(err)
	}
}

//
// Get key from Redis cache
//
func Get(key string, result interface{}) (bool, error) {
	value, err := redisConnection.Get(key).Result()

	// Does not exist
	if err == redis.Nil {
		return false, err
	} else if err != nil { // Error in connection
		services.Fatal(err)
		return false, err
	}

	// UnMarshal Result
	err = json.Unmarshal([]byte(value), result)

	if err != nil {
		services.Fatal(err)
		return false, err
	}

	return true, nil
}

/* End File */
