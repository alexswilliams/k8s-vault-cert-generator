package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func acquireVaultToken(appConfig AppConfig) string {
	var request = buildVaultLoginRequest(appConfig)

	var responseBody = doVaultExchange(request, appConfig)

	var loginResponseRoot VaultLoginRootBlock
	var err = json.Unmarshal(responseBody, &loginResponseRoot)
	panicOnError(err, "Failed to unmarshal vault login response")
	log.Debugf("Response Body after unmarshalling: %+v", loginResponseRoot)

	return loginResponseRoot.Auth.ClientToken
}

func buildVaultLoginRequest(appConfig AppConfig) *http.Request {
	var url = appConfig.VaultAddress + "/v1/auth/kubernetes/" + appConfig.KubernetesClusterName + "/login"
	log.WithField("init.vault-cert.url", url).Infof("New Vault Login request: %s", url)

	payload, err := json.Marshal(map[string]string{"role": appConfig.KubernetesAuthRoleName, "jwt": appConfig.ServiceAccountToken})
	panicOnError(err, "Could not marshal the payload for the vault login request - Role: '%s', JWT: '%s'",
		appConfig.KubernetesAuthRoleName, appConfig.ServiceAccountToken)

	return buildVaultRequest(url, payload)
}

// VaultLoginRootBlock is the response sent back from a login request
type VaultLoginRootBlock struct {
	RequestID     string              `json:"request_id"`
	LeaseID       string              `json:"lease_id"`
	Renewable     bool                `json:"renewable"`
	LeaseDuration int64               `json:"lease_duration"`
	Data          interface{}         `json:"data"`
	WrapInfo      interface{}         `json:"wrap_info"`
	Warnings      interface{}         `json:"warnings"`
	Auth          VaultLoginAuthBlock `json:"auth"`
}

// VaultLoginAuthBlock contains the short-lived token and other metadata from a login request.
type VaultLoginAuthBlock struct {
	ClientToken   string                      `json:"client_token"`
	Accessor      string                      `json:"accessor"`
	Policies      []string                    `json:"policies"`
	TokenPolicies []string                    `json:"token_policies"`
	Metadata      VaultLoginAuthMetadataBlock `json:"metadata"`
	LeaseDuration int64                       `json:"lease_duration"`
	Renewable     bool                        `json:"renewable"`
	EntityID      string                      `json:"entity_id"`
	TokenType     string                      `json:"token_type"`
	Orphan        bool                        `json:"orphan"`
}

// VaultLoginAuthMetadataBlock contains information about the k8s service account that authorised the request.
type VaultLoginAuthMetadataBlock struct {
	Role                     string `json:"role"`
	ServiceAccountName       string `json:"service_account_name"`
	ServiceAccountNamespace  string `json:"service_account_namespace"`
	ServiceAccountSecretName string `json:"service_account_secret_name"`
	ServiceAccountUID        string `json:"service_account_uid"`
}
