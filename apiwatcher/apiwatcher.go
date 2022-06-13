package apiwatcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type DeviceInfo struct {
	Online bool
	Firing bool
	Mode   string
	Name   string
}

type APIInfo struct {
	DevicesInfo map[string]DeviceInfo
	Time        int64
}

type Watcher interface {
	ShowInfo(http.Client) (APIInfo, error)
}

type APIWatcher struct {
	Host string
	Port int
}

type Requester struct {
	Client http.Client
}

func (requester Requester) CallAlarmManager(req *http.Request) (*http.Response, error) {
	response, responseError := requester.Client.Do(req)
	return response, responseError
}

type AlarmManagerRequester interface {
	CallAlarmManager(req *http.Request) (*http.Response, error)
}

type DevicesInfoRequest struct {
	Success bool              `json:"success"`
	Data    map[string]string `json:"data"`
}

type DeviceInfoRequest struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Mode    string `json:"mode"`
	Firing  bool   `json:"firing"`
	Online  bool   `json:"online"`
}

func (watcher APIWatcher) ShowInfo(alarmManager AlarmManagerRequester) (APIInfo, error) {
	var apiInfo APIInfo
	var body []byte

	requestURL := fmt.Sprintf("http://%s:%d/devices", watcher.Host, watcher.Port)
	apiInfo.DevicesInfo = make(map[string]DeviceInfo)

	request, requestError := http.NewRequest("GET", requestURL, bytes.NewReader(body))
	if requestError != nil {
		return apiInfo, requestError
	}
	response, responseErr := alarmManager.CallAlarmManager(request)
	if responseErr != nil {
		return apiInfo, responseErr
	}
	defer response.Body.Close()
	bs, _ := ioutil.ReadAll(response.Body)

	devicesInfo := DevicesInfoRequest{}
	unmarshalErr := json.Unmarshal(bs, &devicesInfo)
	if unmarshalErr != nil {
		return apiInfo, unmarshalErr
	}

	for device_id, device_name := range devicesInfo.Data {
		deviceRequestURL := fmt.Sprintf("http://%s:%d/devices/status/%s", watcher.Host, watcher.Port, device_id)
		var deviceInfo DeviceInfo

		deviceRequest, deviceRequestError := http.NewRequest("GET", deviceRequestURL, bytes.NewReader(body))
		if deviceRequestError != nil {
			return apiInfo, deviceRequestError
		}
		deviceResponse, deviceResponseErr := alarmManager.CallAlarmManager(deviceRequest)
		if deviceResponseErr != nil {
			return apiInfo, responseErr
		}
		defer deviceResponse.Body.Close()
		deviceBodySource, _ := ioutil.ReadAll(deviceResponse.Body)
		deviceInfoRequest := DeviceInfoRequest{}
		unmarshalDeviceErr := json.Unmarshal(deviceBodySource, &deviceInfoRequest)
		if unmarshalDeviceErr != nil {
			return apiInfo, unmarshalDeviceErr
		}
		if deviceInfoRequest.Success == false {
			return apiInfo, errors.New(deviceInfoRequest.Msg)
		}
		deviceInfo.Online = deviceInfoRequest.Online
		deviceInfo.Mode = deviceInfoRequest.Mode
		deviceInfo.Firing = deviceInfoRequest.Firing
		deviceInfo.Name = device_name
		apiInfo.DevicesInfo[device_id] = deviceInfo
	}
	now := time.Now()
	apiInfo.Time = now.Unix()
	return apiInfo, nil
}
