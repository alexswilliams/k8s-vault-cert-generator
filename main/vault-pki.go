package main

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func requestVaultCert(certRequest RequestSpec, appConfig AppConfig, vaultToken string) VaultCertResponseRootBlock {
	var request = buildVaultCertificateRequest(appConfig, certRequest, vaultToken)

	var responseBody = doVaultExchange(request, appConfig)

	var certResponseRoot VaultCertResponseRootBlock
	var err = json.Unmarshal(responseBody, &certResponseRoot)
	panicOnError(err, "Failed to unmarshal the HTTP response body into a certificate response")
	log.Debugf("Response Body after unmarshalling: %+v", certResponseRoot)

	return certResponseRoot
}

func buildVaultCertificateRequest(appConfig AppConfig, certRequest RequestSpec, vaultToken string) *http.Request {
	var url = appConfig.VaultAddress + "/v1/" + certRequest.VaultPath
	log.WithField("init.vault-cert.url", url).Infof("Starting new Vault PKI request: %s", url)

	var payloadMap = make(map[string]string)
	payloadMap["common_name"] = certRequest.CommonName
	payloadMap["format"] = "pem"
	if certRequest.TTL != "" {
		payloadMap["ttl"] = certRequest.TTL
	}
	if len(certRequest.AltNames) > 0 {
		payloadMap["alt_names"] = strings.Join(certRequest.AltNames, ",")
	}

	payload, err := json.Marshal(payloadMap)
	panicOnError(err, "Could not marshal the payload for a vault certificate issue - payload: %+v", payloadMap)

	var request = buildVaultRequest(url, payload)
	request.Header.Set("X-Vault-Token", vaultToken)
	return request
}

// VaultCertResponseRootBlock is a container for the cert data block
type VaultCertResponseRootBlock struct {
	RequestID string                     `json:"request_id"`
	LeaseID   string                     `json:"lease_id"`
	Renewable bool                       `json:"renewable"`
	Data      VaultCertResponseDataBlock `json:"data"`
	WrapInfo  interface{}                `json:"wrap_info"`
	Warnings  interface{}                `json:"warnings"`
	Auth      interface{}                `json:"auth"`
}

// VaultCertResponseDataBlock is the response format from Vault containing certificate data.
type VaultCertResponseDataBlock struct {
	Certificate    string   `json:"certificate"`
	SerialNumber   string   `json:"serial_number"`
	PrivateKey     string   `json:"private_key"`
	PrivateKeyType string   `json:"private_key_type"`
	Expiration     int64    `json:"expiration"`
	IssuingCA      string   `json:"issuing_ca"`
	CAChain        []string `json:"ca_chain"`
}
