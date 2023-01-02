package define

const (
	DEFAULT_COUNTRY             = "CN"                               // Country Name
	DEFAULT_STATE               = "BJ"                               // State Or Province Name
	DEFAULT_LOCALITY            = "HD"                               // Locality Name
	DEFAULT_ORGANIZATION        = "Lab"                              // Organization Name
	DEFAULT_ORGANIZATIONAL_UNIT = "Dev"                              // Organizational Unit Name
	DEFAULT_COMMON_NAME         = "Hello World"                      // Common Name
	DEFAULT_DOMAINS             = "lab.com,*.lab.com,*.data.lab.com" // Domains

	DEFAULT_FOR_K8S = false // Certs For K8S
	DEFAULT_USER    = ""
	DEFAULT_UID     = ""
	DEFAULT_GID     = ""
	DEFAULT_DIR     = "./ssl"
)

var (
	CERT_COUNTRY             string   // Country Name
	CERT_STATE               string   // State Or Province Name
	CERT_LOCALITY            string   // Locality Name
	CERT_ORGANIZATION        string   // Organization Name
	CERT_ORGANIZATIONAL_UNIT string   // Organizational Unit Name
	CERT_COMMON_NAME         string   // Common Name
	CERT_DOMAINS             []string // Domains

	APP_FOR_K8S bool // Certs For K8S
	APP_USER    string
	APP_UID     string
	APP_GID     string
	APP_DIR     string
)
