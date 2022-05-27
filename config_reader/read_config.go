package config

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
