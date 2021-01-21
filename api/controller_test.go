package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-app/app"
	"go-app/mock"
	"go-app/server/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getTestConfig() *config.APIConfig {
	c := config.GetConfigFromFile("test")
	return &c.APIConfig
}

func TestAPI_home(t *testing.T) {

	api := NewTestAPI(getTestConfig())
	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		buildStubs    func(ex *mock.MockExample)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "Success True",
			url:    "/",
			method: http.MethodGet,
			body:   nil,
			buildStubs: func(ex *mock.MockExample) {
				ex.EXPECT().SayHello(gomock.Eq(false)).Times(1).Return("Hello", nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				assert.Equal(t, string(data), "{\"success\":true,\"payload\":\"Hello\"}\n")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ex := mock.NewMockExample(ctrl)

			tt.buildStubs(ex)
			// setting stub
			api.App.Example = ex

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestAPI_saveHello(t *testing.T) {
	helloID, _ := primitive.ObjectIDFromHex("60096eb2f03a83b5ae315c78")
	saveHelloData := &app.SaveHelloOpts{
		Name: "namaste",
	}
	body, _ := json.Marshal(saveHelloData)
	fmt.Println(string(body))
	api := NewTestAPI(getTestConfig())
	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		saveHelloOpts *app.SaveHelloOpts
		buildStubs    func(ex *mock.MockExample)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "When no content is passed",
			url:           "/",
			method:        http.MethodPost,
			saveHelloOpts: saveHelloData,
			buildStubs: func(ex *mock.MockExample) {
				ex.EXPECT().SaveHello(saveHelloData).Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				assert.Equal(t, "{\"error\":[{\"message\":\"Request body must not be empty\",\"type\":\"BadRequest\"}],\"success\":false,\"request_id\":\"\"}\n", string(data))
			},
		},
		{
			name:          "When valid json data is passed",
			url:           "/",
			method:        http.MethodPost,
			body:          bytes.NewReader(body),
			saveHelloOpts: saveHelloData,
			buildStubs: func(ex *mock.MockExample) {
				ex.EXPECT().SaveHello(saveHelloData).Times(1).Return(&app.SaveHelloResp{ID: helloID, Name: "namaste"}, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				assert.Equal(t, "{\"success\":true,\"payload\":{\"id\":\"60096eb2f03a83b5ae315c78\",\"name\":\"namaste\"}}\n", string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ex := mock.NewMockExample(ctrl)
			tt.buildStubs(ex)
			// setting stub
			api.App.Example = ex

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}
