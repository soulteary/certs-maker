package define

const (
	DEFAULT_COUNTRY             = "CN"                               // Country Name
	DEFAULT_STATE               = "BJ"                               // State Or Province Name
	DEFAULT_LOCALITY            = "HD"                               // Locality Name
	DEFAULT_ORGANIZATION        = "Lab"                              // Organization Name
	DEFAULT_ORGANIZATIONAL_UNIT = "Dev"                              // Organizational Unit Name
	DEFAULT_COMMON_NAME         = "Hello World"                      // Common Name
	DEFAULT_DOMAINS             = "lab.com,*.lab.com,*.data.lab.com" // Domains

	DEFAULT_FOR_K8S     = "off" // Certs For K8S
	DEFAULT_FOR_FIREFOX = "off" // Certs For Firefox

	DEFAULT_USER = ""
	DEFAULT_UID  = ""
	DEFAULT_GID  = ""
	DEFAULT_DIR  = "./ssl"

	DEFAULT_EXPIRE_DAYS = "3650"
)

var (
	CERT_COUNTRY             string   // Country Name
	CERT_STATE               string   // State Or Province Name
	CERT_LOCALITY            string   // Locality Name
	CERT_ORGANIZATION        string   // Organization Name
	CERT_ORGANIZATIONAL_UNIT string   // Organizational Unit Name
	CERT_COMMON_NAME         string   // Common Name
	CERT_DOMAINS             []string // Domains

	APP_FOR_K8S     bool // Certs For K8S
	APP_FOR_FIREFOX bool // Certs For Firefox

	APP_USER       string
	APP_UID        string
	APP_GID        string
	APP_OUTPUT_DIR string

	CERT_EXPIRE_DAYS string
)

const DEFAULT_MODE = 0644
