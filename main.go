package main

import (
	"github.com/knadh/koanf/v2"
	"hangoutsb.in/hangout-content-delivery-api/config"
	"hangoutsb.in/hangout-content-delivery-api/logger"
	"hangoutsb.in/hangout-content-delivery-api/router"
	"hangoutsb.in/hangout-content-delivery-api/storage"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var CONFIG = koanf.New(".")

func main() {
	logger.InitLogger()
	config.InitAppConfig(CONFIG)
	logger.SetGlobalLogLevel(CONFIG)
	storage.BlobStorageConnInit(CONFIG)
	router.StartServer(CONFIG)
}
