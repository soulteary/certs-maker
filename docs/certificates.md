# Certificate usage and trust

These instructions apply to v3.8.2 and current source. For known limitations in the v3.8.1 image and binaries, see [release compatibility](installation.md#release-compatibility).

## Certificate modes

### Default mode

The default mode creates a self-signed certificate with the configured SANs, `CA:TRUE`, and server authentication extensions. It produces five files using the first valid domain as the prefix:

```text
lab.com.conf       OpenSSL configuration
lab.com.pem.crt    PEM certificate
lab.com.pem.key    PEM private key
lab.com.der.crt    DER certificate
lab.com.der.key    DER private key
```

A PEM certificate/key pair is suitable for servers that accept PEM input. PEM and DER represent the same material using different encodings; select the format expected by your server, rather than treating them as different trust levels.

### Kubernetes mode

Enable `FOR_K8S=ON` to use `CA:FALSE`, server authentication and encryption extensions, add `*` and `localhost` SAN entries, and append `.k8s` to the default prefix. The five files have names such as `lab.com.k8s.pem.crt` and `lab.com.k8s.pem.key`.

The added `*` entry does not make the certificate valid for every possible hostname. Include the actual DNS names and IP addresses clients use in `CERT_DNS`.

The tool does not create a Secret or configure an Ingress. After generating and verifying the files, create a TLS Secret in the appropriate namespace, for example:

```bash
kubectl create secret tls lab-tls \
  --cert=ssl/lab.com.k8s.pem.crt \
  --key=ssl/lab.com.k8s.pem.key
```

Reference that Secret in your Ingress or server configuration, and configure client trust separately.

### Firefox mode

With a v3.8.2 native binary or a source build, generate the root CA and server certificate:

```bash
./certs-maker "--CERT_DNS=lab.com,127.0.0.1" --FOR_FIREFOX=ON
```

Enable `FOR_FIREFOX=ON` to create a local root CA and sign a separate `CA:FALSE` server certificate containing DNS/IP SANs and server authentication extensions. Along with the ordinary five files, it creates:

```text
lab.com.rootCA.pem    Public root CA certificate to trust
lab.com.rootCA.key    Private root CA key; keep private
lab.com.pem.csr       Server certificate signing request
lab.com.rootCA.srl    CA serial number file
```

Deploy `lab.com.pem.crt` and `lab.com.pem.key` to the server. Import `lab.com.rootCA.pem` into the clients that should trust that server. Do not distribute either private key as a trust certificate.

Kubernetes mode takes precedence when both mode flags are enabled, so that combination does not generate a Firefox root CA.

## Inspect and verify

Inspect dates, subject, issuer, and domain/IP SANs:

```bash
openssl x509 -in ssl/lab.com.pem.crt -noout -subject -issuer -dates -ext subjectAltName
```

For default mode, verify a hostname against the generated self-signed certificate:

```bash
openssl verify -CAfile ssl/lab.com.pem.crt \
  -verify_hostname lab.com ssl/lab.com.pem.crt
```

For Firefox mode, verify the server certificate against its generated root CA:

```bash
openssl verify -CAfile ssl/lab.com.rootCA.pem \
  -verify_hostname lab.com ssl/lab.com.pem.crt
```

If the certificate includes `127.0.0.1`, use `-verify_ip 127.0.0.1` in place of `-verify_hostname lab.com` to check that IP address. A successful result ends with `OK`; a hostname mismatch or missing SAN requires correcting the certificate before deployment.

Check that the certificate and PEM key contain the same public key:

```bash
openssl x509 -in ssl/lab.com.pem.crt -pubkey -noout | openssl sha256
openssl pkey -in ssl/lab.com.pem.key -pubout | openssl sha256
```

The two hashes should match.

## Client trust

Generating a certificate does not install trust into browsers, operating systems, containers, or other clients. Each client must trust the appropriate issuer or explicitly pin the intended certificate, and the hostname/IP it connects to must be covered by the SANs.

For Firefox mode, open Firefox Settings, search for **Certificates**, open the certificate manager, and import the public `lab.com.rootCA.pem` certificate under **Authorities** for website identification. The exact labels may vary by version and language. Restart or reopen the connection if necessary, then verify the served certificate and hostname. See Mozilla's [manual root import instructions](https://wiki.mozilla.org/CA/Changing_Trust_Settings#Trusting_an_Additional_Root_Certificate) and [enterprise CA configuration guide](https://support.mozilla.org/en-US/kb/setting-certificate-authorities-firefox). Other clients may use a separate operating system, application, or container trust store.

For default and Kubernetes self-signed certificates, configure the client's supported trust or certificate-pinning mechanism for that certificate. Kubernetes mode creates a server certificate with `CA:FALSE`; do not present it as a root CA.

For private production environments, control who can sign certificates, protect CA and server private keys, define renewal and replacement procedures, and distribute trust only to the intended clients. The tool does not provide public CA issuance, automatic deployment, a persistent CA service, or automatic renewal.

## Replacement and troubleshooting

- **Untrusted issuer:** import the public issuer certificate into the trust store used by that client, then verify the chain.
- **Hostname mismatch:** inspect SANs and include the actual hostname/IP in `CERT_DNS`. Setting `CERT_CN` alone does not add a SAN.
- **Firefox mode on v3.8.1:** that release omits SANs from the final server certificate; upgrade to v3.8.2 or build current source to include the SANs.
- **Expired or replaced certificate:** generate a replacement, deploy its certificate/key pair, and update client trust when the issuer or pinned certificate changes.
- **Repeated runs:** output with the same prefix is overwritten. In Firefox mode, a repeated run also replaces the root CA key and certificate; clients may need the new public root certificate.

Keep both PEM and DER keys private. See [ownership and key permissions](configuration.md#file-ownership-and-private-keys) before sharing output directories.
