package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	appConfig := getConfig()

	if len(appConfig.CertificateRequestSpecs) == 0 {
		log.Warn("No certificates requested - exiting normally.")
		os.Exit(0)
	}

	var vaultToken = acquireVaultToken(appConfig)

	for _, request := range appConfig.CertificateRequestSpecs {
		certResponse := requestVaultCert(request, appConfig, vaultToken)

		if request.OutputOptions.KubernetesSecretResourceName != "" {
			log.WithFields(log.Fields{
				"init.vault-cert.common-name":   request.CommonName,
				"init.vault-cert.vault-path":    request.VaultPath,
				"init.vault-cert.output-format": request.OutputOptions.Format,
			}).Warn("Kubernetes secret storage is not yet implemented")
			// TODO: Handle posting secrets up to k8s api.
		}
		if request.OutputOptions.DestinationFolderPath != "" {
			writeCertResponseToFiles(request, certResponse)
		}
	}
}

func panicOnError(err error, msg string, v ...interface{}) {
	if err != nil {
		log.Errorf(msg, v...)
		panic(err)
	}
}
