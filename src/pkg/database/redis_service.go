package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
)

func RedisExecuteLuaScript(scriptName string, key string, args ...interface{}) (interface{}, error) {
	keys := []string{}
	if key != "nil" {
		keys = []string{key}
	}

	// Load the Lua script from the external file
	scriptPath := fmt.Sprintf("lua/%s.lua", scriptName)
	luaScript, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		log.Fatalf("Failed to read Lua script: %v", err)
	}

	result, err := GetRedis().Eval(context.Background(), string(luaScript), keys, args).Result()
	if err != nil {
		fmt.Println("Error running Lua script:", err)
		return nil, err
	}

	return result, nil
}

func RedisGeoAdd(key string, latitude float64, longitude float64, reporterUID string) {
	RedisExecuteLuaScript("geoadd", key, longitude, latitude, reporterUID)
}
