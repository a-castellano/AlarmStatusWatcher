package config

import (
	"os"
	"testing"
)

func TestProcessNoConfigFilePresent(t *testing.T) {

	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without any valid config file should fail.")
	} else {
		if err.Error() != "Environment variable ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION is not defined." {
			t.Errorf("Error should be 'Environment variable ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION is not defined.', but error was '%s'.", err.Error())
		}
	}
}

func TestProcessNonExistentConfigFile(t *testing.T) {

	os.Setenv("ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/nonexistent_config_/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without existent config file should fail.")
	}
}

func TestProcessConfigWithoutAnyRequiredField(t *testing.T) {

	os.Setenv("ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_with_no_redis/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without any valid config file should fail.")
	} else {
		if err.Error() != "Fatal error config: no redis field was found." {
			t.Errorf("Error should be 'Fatal error config: no redis field was found.', but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWithNoRedisPort(t *testing.T) {

	os.Setenv("ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_with_no_redis_port/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without reids port should fail.")
	} else {
		if err.Error() != "Fatal error config: no redis port was defined." {
			t.Errorf("Error should be 'Fatal error config: no redis port was defined', but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWithNoAlarmManagerPort(t *testing.T) {

	os.Setenv("ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_with_no_alarmmanager_port/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without alarmmanager port should fail.")
	} else {
		if err.Error() != "Fatal error config: no alarmmanager port was defined." {
			t.Errorf("Error should be 'Fatal error config: no alarmmanager port was defined', but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWithNoNotifyQueue(t *testing.T) {

	os.Setenv("ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_with_no_notify_queue/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without notify queue should fail.")
	} else {
		if err.Error() != "Fatal error config: no notify queue was defined." {
			t.Errorf("Error should be 'Fatal error config: no notify queue was defined', but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWithNoRequiredMail(t *testing.T) {

	os.Setenv("ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_ok_no_required_mail/")
	_, err := ReadConfig()
	if err != nil {
		t.Errorf("ReadConfig method with no required mail config should not fail.")
	}
}

func TestProcessConfigWithoutRequiredMail(t *testing.T) {

	os.Setenv("ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_without_required_mail/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without required mail queue should fail.")
	} else {
		if err.Error() != "Fatal error config: mail config section is required." {
			t.Errorf("Error should be 'Fatal error config: mail config section is required.', but error was '%s'.", err.Error())
		}
	}
}

func TestOkConfig(t *testing.T) {

	os.Setenv("ALARM_STATUS_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_ok/")
	_, err := ReadConfig()
	if err != nil {
		t.Errorf("ReadConfig method with valid config file shouldn't fail. Error was '%s'.", err.Error())
	}
}
