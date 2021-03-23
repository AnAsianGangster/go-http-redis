/*
 * Author: Yang Aobo
 * Telegram: @AnAsianGangster
 * Created At: March 17, 2021
 * Updated At: March 17, 2021
 * Last Modified By: Yang Aobo
 */

/**
 * This package contains all configurations of the main.go server
 *
 *
 * This file creates a singleton redisConfig **config** struct, and expose it to other
 * modules via function **GetRedisConfig**
 */
package config

type config struct {
	ADDRESS  string
	PASSWORD string
	DB       int64
}

var redisConfig *config

func init() {
	//initialize static instance on load
	redisConfig = &config{
		// ADDRESS:  "localhost:6379", // NOTE this is redis local host
		ADDRESS: "redis:6379", // NOTE this is redis docker container
		PASSWORD: "",
		DB:       0,
	}
}

// GetConfig - get redisConfig instance pre-initialized, and expose the
func GetRedisConfig() *config {
	return redisConfig
}
