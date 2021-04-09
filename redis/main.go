package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type RedisInput struct {
	Method   string      `json:"method"`   // set or get
	Addr     string      `json:"address"`  // address of redis to hit
	Password string      `json:"password"` // potentially optional
	DB       int         `json:"db"`       // use default db as 0
	Values   interface{} `json:"values"`   // get use string for key and set use array of strings KEY=VALUE
}

type ReturnGetOutput struct {
	Value string `json:"value"`
}

var code = "com.redis.error"

func RedisHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(RedisInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     obj.Addr,
		Password: obj.Password,
		DB:       obj.DB,
	})

	var data []byte

	switch obj.Method {
	case "get":
		// cast values to get string
		values := obj.Values.(string)
		val, err := rdb.Get(ctx, values).Result()
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
		rgo := &ReturnGetOutput{
			Value: val,
		}
		data, err = json.Marshal(rgo)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	case "set":
		values := obj.Values.([]interface{})
		for _, v := range values {
			split := strings.Split(v.(string), "=")

			err := rdb.Set(ctx, split[0], split[1], 0).Err()
			if err != nil {
				direktivapps.RespondWithError(w, code, err.Error())
				return
			}
		}
	}

	direktivapps.Respond(w, data)
}

func main() {
	direktivapps.StartServer(RedisHandler)
}
