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
