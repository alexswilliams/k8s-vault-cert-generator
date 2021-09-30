package main

import (
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func writeCertResponseToFiles(request RequestSpec, certResponse VaultCertResponseRootBlock) {
	switch request.OutputOptions.Format {
	case PEM:
		writePEMs(request, certResponse)
	case PKCS12:
		writePKCS12(request, certResponse)
	// TODO: Maybe we could do JKS, but given it's been deprecated since JDK9, perhaps not?
	default:
		log.WithFields(log.Fields{
			"init.vault-cert.common-name":   request.CommonName,
			"init.vault-cert.vault-path":    request.VaultPath,
			"init.vault-cert.output-format": request.OutputOptions.Format,
		}).Fatalf("Unknown output format '%s'.", request.OutputOptions.Format)
	}
}

func writePEMs(request RequestSpec, certResponse VaultCertResponseRootBlock) {
	var folder = trimTrailingSlash(request.OutputOptions.DestinationFolderPath)
	var prefix = request.OutputOptions.FileNamePrefix

	writeFile(request, folder, prefix, "certificate.pem", []byte(certResponse.Data.Certificate))
	writeFile(request, folder, prefix, "key.pem", []byte(certResponse.Data.PrivateKey))
	writeFile(request, folder, prefix, "issuing_ca.pem", []byte(certResponse.Data.IssuingCA))
	writeFile(request, folder, prefix, "chain.pem", []byte(strings.Join(certResponse.Data.CAChain, "\n")))
}

func writePKCS12(request RequestSpec, certResponse VaultCertResponseRootBlock) {
	var folder = trimTrailingSlash(request.OutputOptions.DestinationFolderPath)
	var prefix = request.OutputOptions.FileNamePrefix

	pfxBytes, err := createPkcs12(certResponse.Data.PrivateKey, certResponse.Data.Certificate, certResponse.Data.CAChain)
	panicOnError(err, "Failed to create PKCS12 object")
	writeFile(request, folder, prefix, "keystore.p12", pfxBytes)
}

func writeFile(request RequestSpec, folder string, prefix string, suffix string, contents []byte) {
	os.MkdirAll(folder, 0755)
	var path = folder + "/" + prefix + "-" + suffix

	log.WithFields(log.Fields{
		"init.vault-cert.common-name":   request.CommonName,
		"init.vault-cert.vault-path":    request.VaultPath,
		"init.vault-cert.output-format": request.OutputOptions.Format,
		"init.vault-cert.path":          path,
	}).Infof("Writing file: %s", path)
	ioutil.WriteFile(path, []byte(contents), 0644)
}
