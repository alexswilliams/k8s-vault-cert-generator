package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func buildVaultRequest(url string, payload []byte) *http.Request {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	panicOnError(err, "Failed to construct new post request.")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request
}

func doVaultExchange(request *http.Request, appConfig AppConfig) []byte {
	var client = &http.Client{Timeout: time.Second * 5}
	response, err := client.Do(request)
	panicOnError(err, "Failed to exchange POST request for a response.")
	log.WithFields(log.Fields{
		"init.vault-cert.url":         request.RequestURI,
		"init.vault-cert.status-code": response.StatusCode,
	}).Infof("Response: %d (%s)", response.StatusCode, response.Status)

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	panicOnError(err, "Failed to read response body after HTTP call")
	log.Debugf("Headers: %+v; Response Body: %s", response.Header, string(responseBody))

	if response.StatusCode != 200 {
		log.WithFields(log.Fields{
			"init.vault-cert.url":         request.RequestURI,
			"init.vault-cert.status-code": response.StatusCode,
		}).Fatal("Status code was not 200")
	}

	return responseBody
}
