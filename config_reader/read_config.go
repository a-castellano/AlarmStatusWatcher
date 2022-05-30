package config

import (
	"errors"

	viperLib "github.com/spf13/viper"
)

type RabbitmqConfig struct {
	User      string
	Password  string
	Host      string
	Port      int
	QueueName int
}

type RedisServer struct {
	IP       string
	Port     int
	Password string
	Database int
}

type MailServer struct {
	MailFrom     string
	MailDomain   string
	SMTPHost     string
	SMTPPort     int
	SMTPName     string
	SMTPPassword string
}

type NotifyConfig struct {
	NotifyStatusChange    bool
	NotifyOffline         bool
	SendEmailNotification bool
	SendQueueNotification bool
}

type AlarmManager struct {
	Host string
	Port int
}

type Config struct {
	RabbitmqConfig RabbitmqConfig
	RedisServer    RedisServer
	MailServer     MailServer
	NotifyConfig   NotifyConfig
	AlarmManager   AlarmManager
}

func ReadConfig() (Config, error) {

	var configFileLocation string
	var config Config

	var envVariable string = "ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION"

	requiredVariables := []string{"redis", "alarmmanager", "notify"}

	redisRequiredVariables := []string{"ip", "port", "password", "database"}
	//
	//	webServerRequiredVariables := []string{"port"}
	//
	viper := viperLib.New()

	//Look for config file location defined as env var
	viper.BindEnv(envVariable)
	configFileLocation = viper.GetString(envVariable)
	if configFileLocation == "" {
		// Get config file from default location
		return config, errors.New(errors.New("Environment variable ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION is not defined.").Error())
	}

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configFileLocation)

	if err := viper.ReadInConfig(); err != nil {
		return config, errors.New(errors.New("Fatal error reading config file: ").Error() + err.Error())
	}

	for _, requiredVariable := range requiredVariables {
		if !viper.IsSet(requiredVariable) {
			return config, errors.New("Fatal error config: no " + requiredVariable + " field was found.")
		}
	}

	// Redis
	for _, requiredRedisVariable := range redisRequiredVariables {
		if !viper.IsSet("redis." + requiredRedisVariable) {
			return config, errors.New("Fatal error config: no redis " + requiredRedisVariable + " was defined.")
		}

	}
	config.RedisServer.IP = viper.GetString("redis.ip")
	config.RedisServer.Port = viper.GetInt("redis.port")
	config.RedisServer.Password = viper.GetString("redis.password")
	config.RedisServer.Database = viper.GetInt("redis.database")
	return config, nil
}
