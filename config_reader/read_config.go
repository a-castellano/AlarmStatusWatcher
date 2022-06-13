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
	Destination  string
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

	alarmManagerRequiredVariables := []string{"port", "host"}

	notifyRequiredVariables := []string{"online", "statuschange", "queue", "mail"}
	mailRequiredVariables := []string{"mailfrom", "maildomain", "host", "port", "user", "password", "destination"}
	queueRequiredVariables := []string{"host", "port", "user", "password", "queue"}

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

	// AlarmManager
	for _, requiredAlarmManagerVariable := range alarmManagerRequiredVariables {
		if !viper.IsSet("alarmmanager." + requiredAlarmManagerVariable) {
			return config, errors.New("Fatal error config: no alarmmanager " + requiredAlarmManagerVariable + " was defined.")
		}

	}
	config.AlarmManager.Host = viper.GetString("alarmmanager.host")
	config.AlarmManager.Port = viper.GetInt("alarmmanager.port")

	// Notify
	for _, requiredNotifyVariable := range notifyRequiredVariables {
		if !viper.IsSet("notify." + requiredNotifyVariable) {
			return config, errors.New("Fatal error config: no notify " + requiredNotifyVariable + " was defined.")
		}

	}
	config.NotifyConfig.NotifyStatusChange = viper.GetBool("notify.statuschange")
	config.NotifyConfig.NotifyOffline = viper.GetBool("notify.online")
	config.NotifyConfig.SendEmailNotification = viper.GetBool("notify.mail")
	config.NotifyConfig.SendQueueNotification = viper.GetBool("notify.queue")

	// Check if mail is required Mail is required
	if config.NotifyConfig.SendEmailNotification {
		// Mail
		if !viper.IsSet("mail") {
			return config, errors.New("Fatal error config: mail config section is required.")
		} else {
			for _, requiredMailVariable := range mailRequiredVariables {
				if !viper.IsSet("mail." + requiredMailVariable) {
					return config, errors.New("Fatal error config: no mail " + requiredMailVariable + " was defined.")
				}

			}
		}

		config.MailServer.MailFrom = viper.GetString("mail.mailfrom")
		config.MailServer.MailDomain = viper.GetString("mail.maildomain")
		config.MailServer.SMTPHost = viper.GetString("mail.host")
		config.MailServer.SMTPPort = viper.GetInt("mail.port")
		config.MailServer.SMTPName = viper.GetString("mail.user")
		config.MailServer.SMTPPassword = viper.GetString("mail.password")
		config.MailServer.Destination = viper.GetString("mail.destination")

	}

	// Check if queue config is required
	if config.NotifyConfig.SendQueueNotification {
		// Rabbitmq
		if !viper.IsSet("rabbitmq") {
			return config, errors.New("Fatal error config: rabbitmq config section is required.")
		} else {
			for _, requiredQueueVariable := range queueRequiredVariables {
				if !viper.IsSet("rabbitmq." + requiredQueueVariable) {
					return config, errors.New("Fatal error config: no rabbitmq " + requiredQueueVariable + " was defined.")
				}

			}
		}

	}

	return config, nil
}
