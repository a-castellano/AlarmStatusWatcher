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

func TestGetOneDevice(t *testing.T) {

	mock := MockAlarManagerOneDevice{}

	watcher := APIWatcher{Host: "server.local", Port: 8080}
	_, err := watcher.ShowInfo(&mock)

	if err != nil {
		t.Errorf("TestGetOneDevice should not fail, error was '%s'", err.Error())
	}

}
