package main

// RequestSpec is the specification for a single certificate.
type RequestSpec struct {
	VaultPath     string
	CommonName    string
	AltNames      []string
	TTL           string // e.g. "1h", "30m", "60d", etc.
	OutputOptions OutputOptions
}

// OutputOptions gives details about the format and file locations of output files.
// At least one of DestinationFolderPath and KubernetesSecretResourceName must be specified.
type OutputOptions struct {
	Format                       StorageFormat
	FileNamePrefix               string
	DestinationFolderPath        string // For when running as an init container
	KubernetesSecretResourceName string // For when running as a k8s job
}

// StorageFormat enumerates the types of supported certificate store.
type StorageFormat string

const (
	// PEM will generate PEM-encoded text files
	PEM StorageFormat = "PEM"
	// PKCS12 will generate a PKCS12 keystore
	PKCS12 StorageFormat = "PKCS12"
	// JKS isn't supported as it's been deprecated - aspiration for the future maybe.
)
