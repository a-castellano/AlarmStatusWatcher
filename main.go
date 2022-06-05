package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"time"

	apiwatcher "github.com/a-castellano/AlarmStatusWatcher/apiwatcher"
	config_reader "github.com/a-castellano/AlarmStatusWatcher/config_reader"
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

}
