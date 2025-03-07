/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

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

const GENERATE_CMD_TPL = "openssl req -x509 -newkey rsa:2048 -keyout ${file}.pem.key -out ${file}.pem.crt -days ${expire_days} -nodes -config ${file}.conf"
const GENERATE_FILE_PLACEHOLDER = "${file}"
const GENERATE_EXPIRE_DAYS_PLACEHOLDER = "${expire_days}"

const GENERATE_FOR_FF_STEP1 = "openssl genrsa -out ${file}.rootCA.key 2048"
const GENERATE_FOR_FF_STEP2 = "openssl req -utf8 -x509 -new -nodes -key ${file}.rootCA.key -sha256 -days ${expire_days} -out ${file}.rootCA.pem -config ${file}.conf"
const GENERATE_FOR_FF_STEP3 = "openssl genrsa -out ${file}.pem.key 2048"
const GENERATE_FOR_FF_STEP4 = "openssl req -utf8 -new -key ${file}.pem.key -out ${file}.pem.csr -config ${file}.conf"
const GENERATE_FOR_FF_STEP5 = "openssl x509 -req -in ${file}.pem.csr -CA ${file}.rootCA.pem -CAkey ${file}.rootCA.key -CAcreateserial -out ${file}.pem.crt -days ${expire_days} -sha256"

const CONVERT_CRT_TO_DER = "openssl x509 -in ${file}.pem.crt -outform DER -out ${file}.der.crt"
const CONVERT_KEY_TO_DER = "openssl rsa -in ${file}.pem.key -outform DER -out ${file}.der.key"
