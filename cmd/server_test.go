package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/charliemcelfresh/kata/internal/config"
	"github.com/charliemcelfresh/kata/mocks/mock_middlewares"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type Hello struct {
	Hello string `json:"Hello"`
}

func TestRootHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockLogger := mock_middlewares.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.AssignableToTypeOf("string")).Times(1)

	done := make(chan os.Signal)
	go serve(mockLogger, ":8080", done)
	defer func(done chan os.Signal) { done <- os.Interrupt }(done)

	client := &http.Client{}

	testRequest, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		panic(err)
	}

	testRequest.Header.Set("Content-Type", config.Constants["API_KATA_RESPONSE_CONTENT_TYPE"].(string))

	response, err := client.Do(testRequest)
	assert.NoError(t, err)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	var structuredResponse Hello
	err = json.Unmarshal(body, &structuredResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Kata!", structuredResponse.Hello)
}

func TestRootMissingContentType(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockLogger := mock_middlewares.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.AssignableToTypeOf("string")).Times(1)

	done := make(chan os.Signal)
	go serve(mockLogger, ":8080", done)
	defer func(done chan os.Signal) { done <- os.Interrupt }(done)

	response, err := http.Get("http://localhost:8080/")
	assert.NoError(t, err)
	defer response.Body.Close()

	assert.Equal(t, http.StatusUnsupportedMediaType, response.StatusCode)
}
