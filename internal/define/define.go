package define

type CERT struct {
	Country            string
	State              string
	Locality           string
	Organization       string
	OrganizationalUnit string
	CommonName         string
	Domains            []string

	ForK8S  string
	OwnUser string
	OwnUID  string
	OwnGID  string
}

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
