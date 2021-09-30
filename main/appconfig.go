package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// AppConfig is a container for all configuration that enables the app to run - designed to be passed around liberally.
type AppConfig struct {
	ServiceAccountTokenPath    string
	ServiceAccountToken        string
	CertificateRequestJSONPath string
	CertificateRequestSpecs    []RequestSpec
	VaultAddress               string
	KubernetesAuthRoleName     string
	KubernetesClusterName      string
	EnableDebug                bool
	ConsoleFormatter           string
}

// EnvDefaults defines the default values for each environment variable that this program uses.
// Those without defaults, which much always be specified, will panic if these are omitted.
var EnvDefaults = map[string]func() string{
	"REQUEST_SPEC_PATH":          func() string { panic("REQUEST_SPEC_PATH is a required environment variable") },
	"SERVICE_ACCOUNT_TOKEN_PATH": func() string { return "/var/run/secrets/kubernetes.io/serviceaccount/token" },
	"VAULT_ADDR":                 func() string { panic("VAULT_ADDR is a required environment variable") },
	"K8S_AUTH_ROLE_NAME":         func() string { panic("K8S_AUTH_ROLE_NAME is a required environment variable") },
	"KUBERNETES_CLUSTER_NAME":    func() string { panic("KUBERNETES_CLUSTER_NAME is a required environment variable") },
	"ENABLE_DEBUG":               func() string { return "false" },
	"CONSOLE_FORMATTER":          func() string { return "logstash" },
}

func getConfig() (appConfig AppConfig) {
	envVars := buildEnvVars(EnvDefaults)

	appConfig.CertificateRequestJSONPath = envVars["REQUEST_SPEC_PATH"]
	appConfig.ServiceAccountTokenPath = envVars["SERVICE_ACCOUNT_TOKEN_PATH"]
	appConfig.VaultAddress = trimTrailingSlash(envVars["VAULT_ADDR"])
	appConfig.KubernetesAuthRoleName = envVars["K8S_AUTH_ROLE_NAME"]
	appConfig.KubernetesClusterName = envVars["KUBERNETES_CLUSTER_NAME"]
	appConfig.EnableDebug = envVars["ENABLE_DEBUG"] == "true"
	appConfig.ConsoleFormatter = envVars["CONSOLE_FORMATTER"]

	setUpLogger(appConfig)

	log.Debugf(`Using the following configuration:
	• Request Spec Path: %s
	• Service Account Token Path: %s
	• Vault Address: %s
	• Kubernetes Auth Role Name: %s
	• Kubernetes Cluster Name: %s
	• Debug Enabled: %v
	`+"\n",
		appConfig.CertificateRequestJSONPath,
		appConfig.ServiceAccountTokenPath,
		appConfig.VaultAddress,
		appConfig.KubernetesAuthRoleName,
		appConfig.KubernetesClusterName,
		appConfig.EnableDebug)

	appConfig.ServiceAccountToken = readFile(appConfig.ServiceAccountTokenPath)

	var err = json.Unmarshal([]byte(readFile(appConfig.CertificateRequestJSONPath)), &appConfig.CertificateRequestSpecs)
	panicOnError(err, "Failed to unmarshal the user-provided certificate request spec file.")

	return
}

func buildEnvVars(defaults map[string]func() string) (envVars map[string]string) {
	envVars = make(map[string]string)
	for key, defaultValue := range defaults {
		value, found := os.LookupEnv(key)
		if !found || value == "" {
			envVars[key] = defaultValue()
		} else {
			envVars[key] = value
		}
	}
	return
}

func setUpLogger(appConfig AppConfig) {
	log.SetOutput(os.Stdout)
	if appConfig.EnableDebug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	if appConfig.ConsoleFormatter == "logstash" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	}
}

func readFile(path string) string {
	data, err := ioutil.ReadFile(path)
	panicOnError(err, "Failed to read file: '%s'", path)
	return string(data)
}

func trimTrailingSlash(str string) string {
	if strings.HasSuffix(str, "/") {
		return strings.TrimSuffix(str, "/")
	}
	return str
}
