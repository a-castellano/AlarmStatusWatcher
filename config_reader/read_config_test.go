package config

import (
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
