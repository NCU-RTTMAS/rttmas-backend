package database

import (
	"fmt"
	"path/filepath"

	"strings"

	"os"
	"rttmas-backend/pkg/utils/logger"

	"github.com/redis/go-redis/v9"
)

/* LuaScript structure for storing SHA1 Hashes and redis.Script struct reference*/
type LuaScript struct {
	Sha1 string
	*redis.Script
}

/* Global script tables for storing loaded Lua Scripts*/
var LuaScripts map[string]*LuaScript

func LoadLuaScripts(scriptName string, path string) error {
	// Load the Lua script from the external file
	scriptFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read Lua script: %v", err)
	}

	// Cache the luaScript in Redis and store the SHA1 hash
	luaScript := redis.NewScript(string(scriptFile))
	sha1, err := luaScript.Load(ctx, GetRedis()).Result()
	if err != nil {
		return fmt.Errorf("failed to load Lua script into Redis: %v", err)
	}

	// Save the SHA1 hash in the luaScripts map
	LuaScripts[scriptName] = &LuaScript{
		Sha1:   sha1,
		Script: luaScript,
	}
	return nil
}

/* Initialize script loading procecure; Loads all the scripts to Redsis DB */
func InitLuaScripts() error {
	// Initialize the table for Lua scripts
	LuaScripts = make(map[string]*LuaScript)
	luaDir := "lua" // Directory where Lua scripts are stored

	// Use filepath.WalkDir to traverse all files and directories
	err := filepath.WalkDir(luaDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Only load files with .lua extension
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".lua") {
			scriptName := strings.TrimSuffix(d.Name(), ".lua")
			logger.Debug("Script loaded: " + scriptName)

			// Load the Lua script
			err := LoadLuaScripts(scriptName, path)
			if err != nil {
				return fmt.Errorf("failed to load Lua script %s: %v", scriptName, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk through Lua directory: %v", err)
	}

	logger.Info("Lua Scripts loaded")

	return nil
}

/* Caller for Lua Scripts */
func RedisExecuteLuaScript(name string, keys []string, args ...interface{}) (interface{}, error) {
	script, ok := LuaScripts[name]
	if !ok {
		logger.Error(fmt.Sprintf("Lua script with name %s not found", name))
		return nil, fmt.Errorf("lua script with name %s not found", name)
	}
	return script.Run(ctx, GetRedis(), keys, args...).Result()
}

func RedisGeoAdd(key string, latitude float64, longitude float64, reporterUID string) {
	RedisExecuteLuaScript("geoadd", []string{key}, longitude, latitude, reporterUID)
}
