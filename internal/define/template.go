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
