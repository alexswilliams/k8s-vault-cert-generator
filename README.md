# Vault Certificate Issuer Container

This application is designed such that it can be run both as an init container, and as a k8s job.
It takes a JSON specification as its input, and produces certificates as its outputs.

## JSON Specification
The specification is an array of spec objects, which follows roughly the following schema:
```json5
  {
    "VaultPath": "pki-mount/kafka/.....",     // The path to the issuer in Vault
    "CommonName": "My Service Name",          // The cert common name, either a DNS name or a service name
    "AltNames": [                             // Optional list of SANs to add
      "some-test-kafka-san.example.com"
    ],
    "TTL": "1h",                              // Time before expiry, expressed as e.g. "90m", "2h", etc
    "OutputOptions": {                        
      "Format": "PKCS12",                     // Available formats are "PEM" and "PKCS12"                     
      "FileNamePrefix": "kafka-mtls",         // A prefix given to all generated filenames, see below for details.
      "DestinationFolderPath": "/tmp/",       // The folder to place generated files into
      "KubernetesSecretResourceName": "kafka-secrets"  // The name of the k8s resource to save the files into
    }
  }
```
At least one of `DestinationFolderPath` and `KubernetesSecretResourceName` should be specified.
Currently, writing into Kubernetes secrets is not yet supported, and so the resource name key does nothing.

### Formats and Filenames
The `Format` field controls the files that are generated from the vault response.  It has two possible values:
 - `PEM` - will output a set of files in PEM format:
   - `{folder}/{prefix}-certificate.pem` - the client certificate
   - `{folder}/{prefix}-key.pem` - the private key
   - `{folder}/{prefix}-issuing_ca.pem` - the root issuing certificate
   - `{folder}/{prefix}-chain.pem` - all certificates in the chain including the root certificate, but excluding the client certificate.
 - `PKCS12` - will output a single PKCS#12 file:
   - `{folder}/{prefix}-keystore.p12` - contains the client certificate, private key and issuing CAs.

## Keystore Passwords
All keystores generated by this tool are given the hard-coded password `changeit`, as is customary, owing to the poor state of Java Keystore security - see https://neilmadden.blog/2017/11/17/java-keystores-the-gory-details/ for details.

## Environment Variables
The following environment variables are supported for configuring the container:
 - `REQUEST_SPEC_PATH` - the path to a file containing JSON list of the above specification objects.  _Mandatory._
 - `SERVICE_ACCOUNT_TOKEN_PATH` - the path to the kubernetes service account token.  Defaults to `/var/run/secrets/kubernetes.io/serviceaccount/token`.
 - `VAULT_ADDR` - vault address.  _Mandatory._
 - `K8S_AUTH_ROLE_NAME` - the name of the k8s app role in Vault to which policies have been attached.  _Mandatory._ 
 - `KUBERNETES_CLUSTER_NAME` - the name of the kubernetes cluster.  _Mandatory._
 - `ENABLE_DEBUG` - if set to `true` then all requests and response bodies will be output to the console.  Defaults to `false`.
 - `CONSOLE_FORMATTER` - set to `logstash` for easy pickup by the FEKK stack, or anything else for standard console output.  Defaults to `logstash`.

## Building
Running the `build-and-tag.sh` script will compile the application and dockerise it.

The image version is controlled by a parameter within the `build-and-tag.sh` script itself.
Two image tags will always be uploaded:
 - `xxxx:{version}` - e.g. `vault-cert-generator:0.1`
 - `xxxx:{version}-{fromline}` - e.g. `vault-cert-generator:0.1-scratch`

Using the more specific image enables us to update the base build image to address security concerns.  (This no longer applies when using the `FROM scratch` images.)

### Custom CAs
If either vault or k8s require custom certificate authorities, store these in a CA bundle at `main/resources/ca-bundle.pem` and uncomment the line in the dockerfile which includes these.

## Manual Testing
The application can be tested by running it with a suitably configured `config.json` file that exercises all the code paths.
The config file can be located within the `test/resources` folder of the package.

There are two Makefile commands used as part of manual testing, `make run` and `make start-wiremock`.
The wiremock server will respond to vault login requests with a stubbed token response, and to pki requests with a stubbed certificate response.