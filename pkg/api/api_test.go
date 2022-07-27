package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/renanqts/xdpdropper/pkg/logger"
	"github.com/stretchr/testify/assert"
)

type MockXDP struct {
	FakeAddToDrop      func(string) error
	FakeRemoveFromDrop func(string) error
	FakeClose          func()
}

func (m *MockXDP) AddToDrop(strIP string) error {
	return m.FakeAddToDrop(strIP)
}

func (m *MockXDP) RemoveFromDrop(strIP string) error {
	return m.FakeRemoveFromDrop(strIP)
}

func (m *MockXDP) Close() {}

func init() {
	logConfig := logger.NewDefaultConfig()
	logConfig.Level = "debug"
	_ = logger.Init(logConfig)
}

func TestReqUnmarshal(t *testing.T) {
	expectedDrop := drop{IP: "1.2.3.4"}
	dropByte, _ := json.Marshal(expectedDrop)
	expectedFail, _ := json.Marshal([]string{"foo", "bar"})

	testCases := []struct {
		name           string
		payload        []byte
		err            error
		expectedpOuput drop
		statusCode     int
	}{
		{
			name:           "with drop structure",
			payload:        dropByte,
			err:            fmt.Errorf(""),
			expectedpOuput: expectedDrop,
		},
		{
			name:       "with foobar string",
			payload:    expectedFail,
			err:        fmt.Errorf("json: cannot unmarshal array into Go value of type api.drop"),
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("FOOBAR", "/drop", bytes.NewBuffer(tc.payload))
			w := httptest.NewRecorder()
			actualResult, err := reqUnmarshal(w, req)
			if err == nil {
				assert.Equal(t, tc.expectedpOuput, actualResult)
			} else {
				assert.Equal(t, tc.err.Error(), err.Error())
				resp := w.Result()
				assert.Equal(t, tc.statusCode, resp.StatusCode)
			}
		})
	}
}

func TestAddRemove(t *testing.T) {
	apiSuccess := api{
		xdp: &MockXDP{
			FakeAddToDrop: func(strIP string) error {
				return nil
			},
			FakeRemoveFromDrop: func(strIP string) error {
				return nil
			},
		},
	}

	apiFail := api{
		xdp: &MockXDP{
			FakeAddToDrop: func(strIP string) error {
				return fmt.Errorf("foobar")
			},
			FakeRemoveFromDrop: func(strIP string) error {
				return fmt.Errorf("foobar")
			},
		},
	}

	dropTest := drop{
		IP: "1.2.3.4",
	}

	testCases := []struct {
		name       string
		api        api
		payload    drop
		statusCode int
		method     string
	}{
		{
			name:       "add success",
			api:        apiSuccess,
			payload:    dropTest,
			statusCode: http.StatusCreated,
			method:     "POST",
		},
		{
			name:       "add fail",
			api:        apiFail,
			payload:    dropTest,
			statusCode: http.StatusInternalServerError,
			method:     "POST",
		},
		{
			name:       "remove sucess",
			api:        apiSuccess,
			payload:    dropTest,
			statusCode: http.StatusNoContent,
			method:     "DELETE",
		},
		{
			name:       "remove fail",
			api:        apiFail,
			payload:    dropTest,
			statusCode: http.StatusInternalServerError,
			method:     "DELETE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(tc.payload)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, "/drop", bytes.NewBuffer(payload))
			if tc.method == "POST" {
				tc.api.add(w, req)
			} else if tc.method == "DELETE" {
				tc.api.remove(w, req)
			}
			resp := w.Result()
			assert.Equal(t, tc.statusCode, resp.StatusCode)
		})
	}
}
