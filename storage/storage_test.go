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

	newStatus, changedStatusMap, _, _, err := storageInstance.CheckAndUpdate(ctx, devicesInfo)
	if err != nil {
		t.Error("TestNewsReadEmptySet should not fail. Error was ", err.Error())
	}
	if changedStatusMap[key] == "" {
		t.Error("TestNewsReadEmptySet, should not contain empty changedStatusMap variable.")
	}
	if newStatus[key].Name != "Test" {
		t.Error("TestNewsReadEmptySet, name should be Test, not ", newStatus[key].Name)
	}

}

func TestNewsReadNotChanged(t *testing.T) {
	db, mock := redismock.NewClientMock()

	var key string = "ab123"

	expectedValues := make(map[string]string)
	expectedValues["name"] = "Test"
	expectedValues["mode"] = "armed"
	expectedValues["firing"] = "false"
	expectedValues["online"] = "true"

	mock.ExpectHGetAll(key).SetVal(expectedValues)
	storageInstance := Storage{db}
	var ctx = context.TODO()

	deviceInfo := apiwatcher.DeviceInfo{Name: "Test", Mode: "armed", Firing: false, Online: true}
	devicesInfo := make(map[string]apiwatcher.DeviceInfo)

	devicesInfo[key] = deviceInfo

	_, changedStatusMap, _, _, err := storageInstance.CheckAndUpdate(ctx, devicesInfo)
	if err != nil {
		t.Error("TestNewsReadNotChanged should not fail. Error was ", err.Error())
	}
	if changedStatusMap[key] != "" {
		t.Error("TestNewsReadNotChanged, should be empty. It contains ", changedStatusMap[key])
	}

}

func TestNewsReadStartedFirirng(t *testing.T) {
	db, mock := redismock.NewClientMock()

	var key string = "ab123"

	expectedValues := make(map[string]string)
	expectedValues["name"] = "Test"
	expectedValues["mode"] = "armed"
	expectedValues["firing"] = "false"
	expectedValues["online"] = "true"

	mock.ExpectHGetAll(key).SetVal(expectedValues)
	storageInstance := Storage{db}
	var ctx = context.TODO()

	deviceInfo := apiwatcher.DeviceInfo{Name: "Test", Mode: "armed", Firing: true, Online: true}
	devicesInfo := make(map[string]apiwatcher.DeviceInfo)

	devicesInfo[key] = deviceInfo

	_, changedStatusMap, modeChangedMap, onlineChangedMap, err := storageInstance.CheckAndUpdate(ctx, devicesInfo)
	if err != nil {
		t.Error("TestNewsReadStartedFirirng should not fail. Error was ", err.Error())
	}
	if changedStatusMap[key] != "Started Firing" {
		t.Errorf("TestNewsReadStartedFirirng, should be 'Started Firing'. It contains '%s'", changedStatusMap[key])
	}
	if modeChangedMap[key] != true {
		t.Error("TestNewsReadStartedFirirng, modeChangedMap should be true")
	}
	if onlineChangedMap[key] != false {
		t.Error("TestNewsReadStartedFirirng, onlineChangedMap should be false")
	}

}

func TestNewsReadStoppedFirirng(t *testing.T) {
	db, mock := redismock.NewClientMock()

	var key string = "ab123"

	expectedValues := make(map[string]string)
	expectedValues["name"] = "Test"
	expectedValues["mode"] = "armed"
	expectedValues["firing"] = "true"
	expectedValues["online"] = "true"

	mock.ExpectHGetAll(key).SetVal(expectedValues)
	storageInstance := Storage{db}
	var ctx = context.TODO()

	deviceInfo := apiwatcher.DeviceInfo{Name: "Test", Mode: "armed", Firing: false, Online: true}
	devicesInfo := make(map[string]apiwatcher.DeviceInfo)

	devicesInfo[key] = deviceInfo

	_, changedStatusMap, modeChangedMap, onlineChangedMap, err := storageInstance.CheckAndUpdate(ctx, devicesInfo)
	if err != nil {
		t.Error("TestNewsReadStoppedFirirng should not fail. Error was ", err.Error())
	}
	if changedStatusMap[key] != "Stopped Firing" {
		t.Errorf("TestNewsReadStoppedFirirng, should be 'Stopped Firing'. It contains '%s'", changedStatusMap[key])
	}
	if modeChangedMap[key] != true {
		t.Error("TestNewsReadStoppedFirirng, modeChangedMap should be true")
	}
	if onlineChangedMap[key] != false {
		t.Error("TestNewsReadStoppedFirirng, onlineChangedMap should be false")
	}

}

func TestNewsReadBecameOffline(t *testing.T) {
	db, mock := redismock.NewClientMock()

	var key string = "ab123"

	expectedValues := make(map[string]string)
	expectedValues["name"] = "Test"
	expectedValues["mode"] = "armed"
	expectedValues["firing"] = "false"
	expectedValues["online"] = "true"

	mock.ExpectHGetAll(key).SetVal(expectedValues)
	storageInstance := Storage{db}
	var ctx = context.TODO()

	deviceInfo := apiwatcher.DeviceInfo{Name: "Test", Mode: "armed", Firing: false, Online: false}
	devicesInfo := make(map[string]apiwatcher.DeviceInfo)

	devicesInfo[key] = deviceInfo

	_, changedStatusMap, modeChangedMap, onlineChangedMap, err := storageInstance.CheckAndUpdate(ctx, devicesInfo)
	if err != nil {
		t.Error("TestNewsReadBecameOffline should not fail. Error was ", err.Error())
	}
	if changedStatusMap[key] != "Became Offline" {
		t.Errorf("TestNewsReadBecameOffline, should be 'Became Offline'. It contains '%s'", changedStatusMap[key])
	}
	if modeChangedMap[key] != false {
		t.Error("TestNewsReadBecameOffline, modeChangedMap should be false")
	}
	if onlineChangedMap[key] != true {
		t.Error("TestNewsReadBecameOffline, onlineChangedMap should be true")
	}

}

func TestNewsReadBecameOnline(t *testing.T) {
	db, mock := redismock.NewClientMock()

	var key string = "ab123"

	expectedValues := make(map[string]string)
	expectedValues["name"] = "Test"
	expectedValues["mode"] = "armed"
	expectedValues["firing"] = "false"
	expectedValues["online"] = "false"

	mock.ExpectHGetAll(key).SetVal(expectedValues)
	storageInstance := Storage{db}
	var ctx = context.TODO()

	deviceInfo := apiwatcher.DeviceInfo{Name: "Test", Mode: "armed", Firing: false, Online: true}
	devicesInfo := make(map[string]apiwatcher.DeviceInfo)

	devicesInfo[key] = deviceInfo

	_, changedStatusMap, modeChangedMap, onlineChangedMap, err := storageInstance.CheckAndUpdate(ctx, devicesInfo)
	if err != nil {
		t.Error("TestNewsReadBecameOnline should not fail. Error was ", err.Error())
	}
	if changedStatusMap[key] != "Became Online" {
		t.Errorf("TestNewsReadBecameOnline, should be 'Became Online'. It contains '%s'", changedStatusMap[key])
	}
	if modeChangedMap[key] != false {
		t.Error("TestNewsReadBecameOnline, modeChangedMap should be false")
	}
	if onlineChangedMap[key] != true {
		t.Error("TestNewsReadBecameOnline, onlineChangedMap should be true")
	}

}
