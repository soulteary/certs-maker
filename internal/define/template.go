package define

const CERT_BASE_INFO = `
[req]
prompt                  = no
default_bits            = 4096
default_md              = sha256
encrypt_key             = no
string_mask             = utf8only

distinguished_name      = cert_distinguished_name
req_extensions          = req_x509v3_extensions
x509_extensions         = req_x509v3_extensions	
`

const CERT_EXTENSIONS = `
[req_x509v3_extensions]
basicConstraints        = critical,CA:true
subjectKeyIdentifier    = hash
keyUsage                = critical,digitalSignature,keyCertSign,cRLSign
extendedKeyUsage        = critical,serverAuth
subjectAltName          = @alt_names
`

const CERT_EXTENSIONS_K8S = `
[req_x509v3_extensions]
basicConstraints = CA:FALSE
nsCertType = server
nsComment = "OpenSSL Generated Server Certificate"
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer:always
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
`

const GENERATE_CMD_TPL = "openssl req -x509 -newkey rsa:2048 -keyout ${file}.key -out ${file}.crt -days 3650 -nodes -config ${file}.conf"
const GENERATE_CMD_PLACEHOLDER = "${file}"

const GENERATE_FOR_FF_STEP1 = "openssl genrsa -out ${file}.rootCA.key 2048"
const GENERATE_FOR_FF_STEP2 = "openssl req -x509 -new -nodes -key ${file}.rootCA.key -sha256 -days 3650 -out ${file}.rootCA.pem -config ${file}.conf"
const GENERATE_FOR_FF_STEP3 = "openssl genrsa -out ${file}.key 2048"
const GENERATE_FOR_FF_STEP4 = "openssl req -new -key ${file}.key -out ${file}.csr -config ${file}.conf"
const GENERATE_FOR_FF_STEP5 = "openssl x509 -req -in ${file}.csr -CA ${file}.rootCA.pem -CAkey ${file}.rootCA.key -CAcreateserial -out ${file}.crt -days 3650 -sha256"
