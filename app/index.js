const dotenv = require("dotenv");
const { writeFileSync, existsSync, mkdirSync } = require("fs");

const defaults = `
# Country Name
CERT_C=CN
#State Or Province Name
CERT_ST=BJ
# Locality Name
CERT_L=HD
# Organization Name
CERT_O=Lab
# Organizational Unit Name
CERT_OU=Dev
# Common Name
CERT_CN=Hello World
# Domians
CERT_DNS=lab.com;*.lab.com;*.data.lab.com
`;

function parseConfigFromString(str) {
  return dotenv.parse(Buffer.from(str), { debug: false });
}

function parseUserInputs() {
  return Object.keys(process.env)
    .filter((n) => n.startsWith("CERT_"))
    .reduce((prev, label) => {
      prev[label] = process.env[label];
      return prev;
    }, {});
}

function error(msg) {
  console.error(msg);
  process.exit(1);
}

function merge(userInputs, defaults) {
  let config = Object.assign({}, defaults);
  Object.keys(userInputs).forEach((key) => {
    let val = userInputs[key];
    if (val) val = (val + "").trim();
    switch (key.toUpperCase()) {
      case "CERT_C":
        val = userInputs[key];
        if (!val.match(/^(\w){2}$/)) return error(`${key} invaild.`);
        config[key] = val.toUpperCase();
        break;
      case "CERT_ST":
      case "CERT_L":
        if (val) config[key] = val.toUpperCase();
        break;
      case "CERT_O":
      case "CERT_OU":
      case "CERT_CN":
        if (val) config[key] = userInputs[key];
        break;
      case "CERT_DNS":
        if (val.match(/([\.\w\*]+(\;)?){1,}/)) config[key] = userInputs[key];
        break;
    }
  });

  config['CERT_DNS'] = config['CERT_DNS'].split(";").filter((domain) => domain.match(/^[\*\w]+(\.)?.+\w$/));
  return config;
}

function generateConfFile(c) {
  const dns = c.CERT_DNS.map((dns, idx) => `DNS.${idx + 1} = ${dns}`).join("\n");

  const tplCerts = `
[req]
prompt                  = no
default_bits            = 4096
default_md              = sha256
encrypt_key             = no
string_mask             = utf8only

distinguished_name      = cert_distinguished_name
req_extensions          = req_x509v3_extensions
x509_extensions         = req_x509v3_extensions

[cert_distinguished_name]
C  = ${c.CERT_C}
ST = ${c.CERT_ST}
L  = ${c.CERT_L}
O  = ${c.CERT_O}
OU = ${c.CERT_OU}
CN = ${c.CERT_CN}

[req_x509v3_extensions]
basicConstraints        = critical,CA:true
subjectKeyIdentifier    = hash
keyUsage                = critical,digitalSignature,keyCertSign,cRLSign
extendedKeyUsage        = critical,serverAuth
subjectAltName          = @alt_names

[alt_names]
${dns}
`;

  return tplCerts;
}

const inputs = parseUserInputs();
console.log(`User Input:`, inputs);

const config = merge(inputs, parseConfigFromString(defaults));

if (!existsSync("./ssl")) mkdirSync("./ssl");

const fileName = config.CERT_DNS[0];
writeFileSync(`./ssl/${fileName}.conf`, generateConfFile(config));
writeFileSync(`./generate.sh`,`openssl req -x509 -newkey rsa:2048 -keyout ssl/${fileName}.key -out ssl/${fileName}.crt -days 3600 -nodes -config ssl/${fileName}.conf`);
