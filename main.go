package main

import (
	"context"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"time"

	apiwatcher "github.com/a-castellano/AlarmStatusWatcher/apiwatcher"
	config_reader "github.com/a-castellano/AlarmStatusWatcher/config_reader"
	storage "github.com/a-castellano/AlarmStatusWatcher/storage"
	goredis "github.com/go-redis/redis/v8"
)

func main() {
	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	alarmManagerRequester := apiwatcher.Requester{Client: client}
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "AlarmManager")
	if e == nil {
		log.SetOutput(logwriter)
		// Remove date prefix
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}

	config, errConfig := config_reader.ReadConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
		return
	}

	watcher := apiwatcher.APIWatcher{Host: config.AlarmManager.Host, Port: config.AlarmManager.Port}

	apiInfo, err := watcher.ShowInfo(alarmManagerRequester)
	if err != nil {

		log.Fatal(err)
		return
	}
	fmt.Println(apiInfo)
	redisAddress := fmt.Sprintf("%s:%d", config.RedisServer.IP, config.RedisServer.Port)

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     redisAddress,
		Password: config.RedisServer.Password,
		DB:       config.RedisServer.Database,
	})
	ctx := context.Background()
	storageInstance := storage.Storage{RedisClient: rdb}
	newStatusMap, changedStatusMap, _ := storageInstance.CheckAndUpdate(ctx, apiInfo.DevicesInfo)
	fmt.Println(newStatusMap)
	fmt.Println(changedStatusMap)
	apiInfo.DevicesInfo = newStatusMap
	newStatusMap2, changedStatusMap2, _ := storageInstance.CheckAndUpdate(ctx, apiInfo.DevicesInfo)
	fmt.Println(newStatusMap2)
	fmt.Println(changedStatusMap2)
}
