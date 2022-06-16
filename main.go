package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"log/syslog"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"time"

	apiwatcher "github.com/a-castellano/AlarmStatusWatcher/apiwatcher"
	config_reader "github.com/a-castellano/AlarmStatusWatcher/config_reader"
	storage "github.com/a-castellano/AlarmStatusWatcher/storage"
	goredis "github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

func sendMessageByQueue(rabbitmqConfig config_reader.RabbitmqConfig, messageToSend string) error {

	dialString := fmt.Sprintf("amqp://%s:%s@%s:%d/", rabbitmqConfig.User, rabbitmqConfig.Password, rabbitmqConfig.Host, rabbitmqConfig.Port)

	conn, errDial := amqp.Dial(dialString)
	defer conn.Close()

	if errDial != nil {
		return errDial
	}

	channel, errChannel := conn.Channel()
	defer channel.Close()
	if errChannel != nil {
		return errChannel
	}

	queue, errQueue := channel.QueueDeclare(
		rabbitmqConfig.QueueName, // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if errQueue != nil {
		return errQueue
	}

	// send Job

	err := channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(messageToSend),
		})

	if err != nil {
		return err
	}
	return nil

}

func sendEmail(config config_reader.Config, messageToSend string) {

	fromMail := fmt.Sprintf("%s@%s", config.MailServer.MailFrom, config.MailServer.MailDomain)
	from := mail.Address{"", fromMail}
	to := mail.Address{"", config.MailServer.Destination}
	subj := "Alarm Status Changed"

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	var message string
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + messageToSend

	// Connect to the SMTP Server
	servername := fmt.Sprintf("%s:%d", config.MailServer.SMTPHost, config.MailServer.SMTPPort)

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", config.MailServer.SMTPName, config.MailServer.SMTPPassword, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
	log.Println("Mail sent successfully")
}

func checkStatus(ctx context.Context, config config_reader.Config, storageInstance storage.Storage, alarmManagerRequester apiwatcher.Requester) {

	watcher := apiwatcher.APIWatcher{Host: config.AlarmManager.Host, Port: config.AlarmManager.Port}

	for range time.Tick(time.Second * 1) {
		log.Println("Checking api status.")
		apiInfo, apiInfoErr := watcher.ShowInfo(alarmManagerRequester)
		if apiInfoErr != nil {

			log.Fatal(apiInfoErr)
			return
		}
		newStatusMap, changedStatusMap, modeChangedMap, onlineChangedMap, checkAndUpdateErr := storageInstance.CheckAndUpdate(ctx, apiInfo.DevicesInfo)
		if checkAndUpdateErr != nil {
			log.Fatal(checkAndUpdateErr)
			return
		}
		apiInfo.DevicesInfo = newStatusMap
		for deviceID, message := range changedStatusMap {
			if len(message) > 0 {
				notificationMessage := fmt.Sprintf("%s - %s", apiInfo.DevicesInfo[deviceID].Name, message)
				if config.NotifyConfig.SendEmailNotification {
					if (config.NotifyConfig.NotifyOffline == true && onlineChangedMap[deviceID] == true) || (config.NotifyConfig.NotifyStatusChange == true && modeChangedMap[deviceID] == true) {
						sendEmail(config, notificationMessage)
					}
				}
				if config.NotifyConfig.SendQueueNotification {
					if (config.NotifyConfig.NotifyOffline == true && onlineChangedMap[deviceID] == true) || (config.NotifyConfig.NotifyStatusChange == true && modeChangedMap[deviceID] == true) {
						sendError := sendMessageByQueue(config.RabbitmqConfig, notificationMessage)
						if sendError != nil {
							log.Fatal(sendError)
						}
					}
				}
			}

		}
	}
}

func main() {
	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	alarmManagerRequester := apiwatcher.Requester{Client: client}
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "AlarmStatusWatcher")
	if e == nil {
		log.SetOutput(logwriter)
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}

	config, errConfig := config_reader.ReadConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
		return
	}

	redisAddress := fmt.Sprintf("%s:%d", config.RedisServer.IP, config.RedisServer.Port)

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     redisAddress,
		Password: config.RedisServer.Password,
		DB:       config.RedisServer.Database,
	})

	ctx := context.Background()

	redisErr := redisClient.Set(ctx, "checkKey", "key", 1000000).Err()
	if redisErr != nil {
		panic(redisErr)
	}
	storageInstance := storage.Storage{RedisClient: redisClient}

	checkStatus(ctx, config, storageInstance, alarmManagerRequester)

}
