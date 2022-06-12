package storage

import (
	"context"
	"testing"

	apiwatcher "github.com/a-castellano/AlarmStatusWatcher/apiwatcher"
	redismock "github.com/go-redis/redismock/v8"
)

func TestNewsReadEmptySet(t *testing.T) {
	db, mock := redismock.NewClientMock()

	var key string = "ab123"
	mock.ExpectHGetAll(key).RedisNil()

	storageInstance := Storage{db}
	var ctx = context.TODO()

	deviceInfo := apiwatcher.DeviceInfo{Name: "Test", Mode: "test", Firing: false, Online: true}
	devicesInfo := make(map[string]apiwatcher.DeviceInfo)

	devicesInfo[key] = deviceInfo

	_, changedStatusMap, err := storageInstance.CheckAndUpdate(ctx, devicesInfo)
	if err != nil {
		t.Error("TestNewsReadEmptySet should not fail. Error was ", err.Error())
	}
	if changedStatusMap[key] == "" {
		t.Error("TestNewsReadEmptySet, should not contain empty changedStatusMap variable.")
	}

}
