package config

type config struct {
	ADDRESS  string
	PASSWORD string
	DB       int64
}

var singleton *config

func init() {
	//initialize static instance on load
	singleton = &config{
		ADDRESS:  "localhost:6379",
		PASSWORD: "",
		DB:       0,
	}
}

//GetConfig - get singleton instance pre-initialized
func GetRedisConfig() *config {
	return singleton
}
