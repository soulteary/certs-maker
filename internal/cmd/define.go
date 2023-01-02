package cmd

const (
	DEFAULT_COUNTRY             = "CN"                               // Country Name
	DEFAULT_STATE               = "BJ"                               // State Or Province Name
	DEFAULT_LOCALITY            = "HD"                               // Locality Name
	DEFAULT_ORGANIZATION        = "Lab"                              // Organization Name
	DEFAULT_ORGANIZATIONAL_UNIT = "Dev"                              // Organizational Unit Name
	DEFAULT_COMMON_NAME         = "Hello World"                      // Common Name
	DEFAULT_DOMAINS             = "lab.com,*.lab.com,*.data.lab.com" // Domains

	DEFAULT_FORK8S = "OFF" // Certs For K8S
)

const (
	ENV_KEY_COUNTRY           = "CERT_C"
	ENV_KEY_STATE             = "CERT_ST"
	ENV_KEY_LOCALITY          = "CERT_L"
	ENV_KEY_ORGANIZATION      = "CERT_O"
	ENV_KEY_ORGANIZATION_UNIT = "CERT_OU"
	ENV_KEY_COMMON_NAME       = "CERT_CN"
	ENV_KEY_DOMAINS           = "CERT_DNS"

	ENV_KEY_FOR_K8S = "FOR_K8S"
	ENV_KEY_USER    = "USER"
	ENV_KEY_UID     = "UID"
	ENV_KEY_GID     = "GID"
)
