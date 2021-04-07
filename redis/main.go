package main

import (
	"context"
	"encoding/json"
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

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.redis.error",
		ErrorMessage: "",
	}

	obj := new(RedisInput)
	direktivapps.ReadIn(obj, g)

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
			g.ErrorMessage = err.Error()
			direktivapps.WriteError(g)
		}
		rgo := &ReturnGetOutput{
			Value: val,
		}
		data, err = json.Marshal(rgo)
		if err != nil {
			g.ErrorMessage = err.Error()
			direktivapps.WriteError(g)
		}
	case "set":
		values := obj.Values.([]interface{})
		for _, v := range values {
			split := strings.Split(v.(string), "=")

			err := rdb.Set(ctx, split[0], split[1], 0).Err()
			if err != nil {
				g.ErrorMessage = err.Error()
				direktivapps.WriteError(g)
			}
		}
	}

	direktivapps.WriteOut(data, g)
}
