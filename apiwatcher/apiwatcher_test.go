package apiwatcher

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

type MockAlarManagerOneDevice struct {
	CallCounter int
}

func (m *MockAlarManagerOneDevice) CallAlarmManager(req *http.Request) (*http.Response, error) {
	var client http.Client
	if m.CallCounter == 0 {
		client = http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"success":true,"data":{"deviceid":"Home Alarm"}}`))}}}
	} else {
		client = http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"success":true,"msg":"","mode":"disarmed","firing":false,"online":true}`))}}}
	}

	response, responseError := client.Do(req)
	m.CallCounter++
	return response, responseError
}

type MockAlarManagerErrorFirstRequest struct {
	CallCounter int
}

func (m *MockAlarManagerErrorFirstRequest) CallAlarmManager(req *http.Request) (*http.Response, error) {
	var client http.Client
	if m.CallCounter == 0 {
		client = http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`"success":true,"data":{"deviceid":"Home Alarm"}}`))}}}
	} else {
		client = http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"success":true,"msg":"","mode":"disarmed","firing":false,"online":true}`))}}}
	}

	response, responseError := client.Do(req)
	m.CallCounter++
	return response, responseError
}

type MockAlarManagerErrorSecondRequest struct {
	CallCounter int
}

func (m *MockAlarManagerErrorSecondRequest) CallAlarmManager(req *http.Request) (*http.Response, error) {
	var client http.Client
	if m.CallCounter == 0 {
		client = http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"success":true,"data":{"deviceid":"Home Alarm"}}`))}}}
	} else {
		client = http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`"success":true,"msg":"","mode":"disarmed","firing":false,"online":true}`))}}}
	}

	response, responseError := client.Do(req)
	m.CallCounter++
	return response, responseError
}

type MockAlarManagerSecondRequestWithError struct {
	CallCounter int
}

func (m *MockAlarManagerSecondRequestWithError) CallAlarmManager(req *http.Request) (*http.Response, error) {
	var client http.Client
	if m.CallCounter == 0 {
		client = http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"success":true,"data":{"deviceid":"Home Alarm"}}`))}}}
	} else {
		client = http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"success":false,"msg":"Failed","mode":"disarmed","firing":false,"online":true}`))}}}
	}

	response, responseError := client.Do(req)
	m.CallCounter++
	return response, responseError
}

func TestGetOneDevice(t *testing.T) {

	mock := MockAlarManagerOneDevice{}

	watcher := APIWatcher{Host: "server.local", Port: 8080}
	apiInfo, err := watcher.ShowInfo(&mock)

	if err != nil {
		t.Errorf("TestGetOneDevice should not fail, error was '%s'", err.Error())
	}

	if len(apiInfo.DevicesInfo) != 1 {
		t.Errorf("API info should have only one device infor not, %d.", len(apiInfo.DevicesInfo))
	}

}

func TestGetOneDeviceErrorOnList(t *testing.T) {

	mock := MockAlarManagerErrorFirstRequest{}

	watcher := APIWatcher{Host: "server.local", Port: 8080}
	_, err := watcher.ShowInfo(&mock)

	if err == nil {
		t.Errorf("TestGetOneDeviceErrorOnList should fail.")
	}

}

func TestGetOneDeviceErrorOnDeviceInfo(t *testing.T) {

	mock := MockAlarManagerErrorSecondRequest{}

	watcher := APIWatcher{Host: "server.local", Port: 8080}
	_, err := watcher.ShowInfo(&mock)

	if err == nil {
		t.Errorf("TestGetOneDeviceErrorOnDeviceInfo should fail.")
	}

}

func TestGetOneDeviceErrorOnDeviceRequest(t *testing.T) {

	mock := MockAlarManagerSecondRequestWithError{}

	watcher := APIWatcher{Host: "server.local", Port: 8080}
	_, err := watcher.ShowInfo(&mock)

	if err == nil {
		t.Errorf("TestGetOneDeviceErrorOnDeviceInfo should fail.")
	}

}

func TestRequester(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"success":true,"data":{"deviceid":"Home Alarm"}}`))}}}

	requester := Requester{Client: client}
	request, _ := http.NewRequest("GET", "http://test.local/api", nil)
	_, err := requester.CallAlarmManager(request)

	if err != nil {
		t.Errorf("TestRequester shouldn't' fail.")
	}

}
