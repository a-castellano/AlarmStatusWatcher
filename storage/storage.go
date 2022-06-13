package storage

import (
	"context"
	"fmt"
	"strings"

	apiwatcher "github.com/a-castellano/AlarmStatusWatcher/apiwatcher"
	goredis "github.com/go-redis/redis/v8"
)

type AlarmStatus struct {
	Online bool   `redis:"online"`
	Firing bool   `redis:"firing"`
	Mode   string `redis:"mode"`
	Name   string `redis:"name"`
}

type Storage struct {
	RedisClient *goredis.Client
}

func (storage Storage) CheckAndUpdate(ctx context.Context, devicesInfo map[string]apiwatcher.DeviceInfo) (map[string]apiwatcher.DeviceInfo, map[string]string, error) {
	newStatusMap := make(map[string]apiwatcher.DeviceInfo)
	changedStatusMap := make(map[string]string)
	for deviceId, newDeviceInfo := range devicesInfo {
		var storedAlarmStatus AlarmStatus
		storedAlarmStatusError := storage.RedisClient.HGetAll(ctx, deviceId).Scan(&storedAlarmStatus)
		if storedAlarmStatusError != goredis.Nil {
			if storedAlarmStatusError != nil {
				fmt.Println("ERRROR", storedAlarmStatusError)
				return newStatusMap, changedStatusMap, storedAlarmStatusError
			}
		} else { // Value has not been set yet
			storedAlarmStatus.Name = ""
			storedAlarmStatus.Mode = "Not Set"
			storedAlarmStatus.Firing = false
			storedAlarmStatus.Online = false
		}

		// Compare Values
		changedStatusMap[deviceId] = ""
		if storedAlarmStatus.Name != newDeviceInfo.Name {
			changedStatusMap[deviceId] = fmt.Sprintf("%sChanged Name to %s ", changedStatusMap[deviceId], newDeviceInfo.Name)
		}
		storedAlarmStatus.Name = newDeviceInfo.Name
		if storedAlarmStatus.Mode != newDeviceInfo.Mode {
			changedStatusMap[deviceId] = fmt.Sprintf("%sChanged Mode from %s to %s ", changedStatusMap[deviceId], storedAlarmStatus.Mode, newDeviceInfo.Mode)
		}
		storedAlarmStatus.Mode = newDeviceInfo.Mode
		if storedAlarmStatus.Firing != newDeviceInfo.Firing {
			if newDeviceInfo.Firing == true {
				changedStatusMap[deviceId] = fmt.Sprintf("%sStarted Firing ", changedStatusMap[deviceId])
			} else {
				changedStatusMap[deviceId] = fmt.Sprintf("%sStopped Firing ", changedStatusMap[deviceId])
			}
		}
		storedAlarmStatus.Firing = newDeviceInfo.Firing
		if storedAlarmStatus.Online != newDeviceInfo.Online {
			if newDeviceInfo.Online == true {
				changedStatusMap[deviceId] = fmt.Sprintf("%sBecame Online ", changedStatusMap[deviceId])
			} else {
				changedStatusMap[deviceId] = fmt.Sprintf("%sBecame Offline ", changedStatusMap[deviceId])
			}
		}
		storedAlarmStatus.Online = newDeviceInfo.Online

		storage.RedisClient.HSet(ctx, deviceId, "name", newDeviceInfo.Name)
		storage.RedisClient.HSet(ctx, deviceId, "mode", newDeviceInfo.Mode)
		storage.RedisClient.HSet(ctx, deviceId, "online", newDeviceInfo.Online)
		storage.RedisClient.HSet(ctx, deviceId, "firing", newDeviceInfo.Firing)

		newStatusMap[deviceId] = newDeviceInfo
		changedStatusMap[deviceId] = strings.TrimSpace(changedStatusMap[deviceId])
	}
	return newStatusMap, changedStatusMap, nil
}