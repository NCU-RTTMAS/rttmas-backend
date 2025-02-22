package redis

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	"embed"
	// "os"
	// "rttmas-backend/utils"
	"rttmas-backend/utils/logger"
	"strings"

	"github.com/redis/go-redis/v9"
)

/* LuaScript structure for storing SHA1 Hashes and redis.Script struct reference*/
type LuaScript struct {
	Sha1 string
	*redis.Script
}

//go:embed lua/*
var LuaScriptFS embed.FS

/* Global script tables for storing loaded Lua Scripts*/
var LuaScripts map[string]*LuaScript

func LoadLuaScripts(scriptName string, path string) error {
	// Load the Lua script from the embedded filesystem
	scriptFile, err := LuaScriptFS.ReadFile(path)
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

	// Use LuaScriptFS to traverse all embedded files
	err := fs.WalkDir(LuaScriptFS, "lua", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Only load files with .lua extension
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".lua") {
			scriptName := strings.TrimSuffix(filepath.Base(d.Name()), ".lua")
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
func DeleteKeysWithPrefix(prefix string) error {
	var cursor uint64
	for {
		// SCAN command to find keys with the specified prefix
		keys, nextCursor, err := GetRedis().Scan(context.Background(), cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		// Delete keys if any are found
		if len(keys) > 0 {
			if _, err := GetRedis().Del(context.Background(), keys...).Result(); err != nil {
				return err
			}
		}

		// Check if we've iterated through all keys
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}
