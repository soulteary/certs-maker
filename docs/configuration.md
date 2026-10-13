# Configuration reference

Every option supports an environment variable and a matching `--NAME=value` CLI flag. The following defaults and precedence apply to v3.8.2 and current source; see [release compatibility](installation.md#release-compatibility) for v3.8.1 limitations.

| Variable / CLI name | Default | Meaning |
| --- | --- | --- |
| `CERT_C` | `CN` | Country code; use a standard two-letter code. See validation below. |
| `CERT_ST` | `BJ` | State or province. |
| `CERT_L` | `HD` | Locality; converted to uppercase. |
| `CERT_O` | `Lab` | Organization. |
| `CERT_OU` | `Dev` | Organizational unit. |
| `CERT_CN` | `Hello World` | Subject common name; hostnames belong in `CERT_DNS` SANs. |
| `CERT_DNS` | `lab.com,*.lab.com,*.data.lab.com` | Comma-separated DNS names, wildcard domains, IPv4, or IPv6 addresses. |
| `FOR_K8S` | `off` | Select Kubernetes certificate extensions and naming. |
| `FOR_FIREFOX` | `off` | Create a local root CA and signed server certificate. |
| `USER` | empty | Name used by the container ownership helper. Requires both `UID` and `GID`. |
| `UID` | empty | Numeric owner UID. Requires both `USER` and `GID`. |
| `GID` | empty | Numeric owner GID. Requires both `USER` and `UID`. |
| `DIR` | `./ssl` | Output directory in the process's filesystem; see path limitations below. |
| `CUSTOM_FILE_NAME` | empty | Override the output base name, without an extension. |
| `EXPIRE_DAYS` | `3650` | Certificate validity in days; integer from `1` to `2147483647`. |

Environment example:

```bash
CERT_DNS="lab.com,127.0.0.1" EXPIRE_DAYS=30 ./certs-maker
```

Equivalent CLI example:

```bash
./certs-maker "--CERT_DNS=lab.com,127.0.0.1" --EXPIRE_DAYS=30
```

Quote an entire flag argument when its value contains spaces, for example `"--CERT_CN=Lab Server"`. Give each flag its own argument; do not combine `--CERT_DNS=... --FOR_K8S=ON` inside one pair of quotes.

## Precedence and accepted values

A non-empty, explicitly supplied CLI value overrides the environment variable, which overrides the built-in default. Omitting a CLI flag leaves the environment variable in control. Empty CLI or environment values fall back to the remaining source/default.

This includes explicitly choosing a default-valued option: `--CERT_C=CN` overrides `CERT_C=US`, and `--FOR_K8S=OFF` disables an environment setting of `FOR_K8S=ON`.

For mode flags, `on`, `true`, and `1` are true, ignoring case. All other values are false; use `ON` and `OFF` to keep configurations easy to read. If both flags are true, Kubernetes mode takes precedence and no Firefox root CA is created.

`CERT_DNS` trims whitespace, drops invalid entries, removes duplicates, and lowercases DNS names. Supply at least one valid entry. Bare names such as `localhost` are not accepted by the normal input parser; Kubernetes mode adds `localhost` separately. A wildcard such as `*.lab.com` does not replace the explicit `lab.com` name or cover arbitrary deeper subdomains.

Country validation checks only for two ASCII letters, digits, or underscores, then uppercases the value. It does not check ISO country assignments; other formats fall back to `CN`. Use a standard two-letter country code.

`EXPIRE_DAYS` must be an integer from `1` to `2147483647`. Invalid values fail generation before any output is written. Choose a validity period suitable for your deployment and renewal policy; the 3650-day default is not a recommendation for every environment.

## Output paths and file names

The default prefix comes from the first valid `CERT_DNS` entry. Kubernetes mode appends `.k8s`; a non-empty `CUSTOM_FILE_NAME` replaces the full prefix, including that suffix. Use a simple base name such as `lab-server`, without spaces, path separators, or shell characters.

For example, the following writes `ssl/lab-server.pem.crt`, `ssl/lab-server.pem.key`, the DER pair, and the OpenSSL configuration:

```bash
./certs-maker "--CERT_DNS=lab.com" --DIR=ssl --CUSTOM_FILE_NAME=lab-server
```

`DIR` is normalized rather than preserved verbatim: custom paths are lowercased, `..` and `./` sequences are removed, unsupported characters are stripped, and leading/trailing slashes are removed. Use simple relative paths such as `ssl` or `certs/server`. Absolute paths, spaces, and case-sensitive directory names may not resolve as you expect.

With Docker, `DIR` refers to a container path, not a host path. To write into another host directory, change the source of the bind mount and keep its destination `/ssl`, for example `-v "$(pwd)/certificates:/ssl"`. Changing `DIR` without a corresponding mount can leave output inside the temporary container.

Rerunning with the same output directory and prefix overwrites existing files and generates new private keys. The tool has no reuse or renewal mode. Move existing output aside or choose another directory/prefix before generating a replacement.

## File ownership and private keys

The ownership helper only runs when `USER`, `UID`, and `GID` are all non-empty. It uses Alpine-compatible account tools and requires permission to create an account/group and change file ownership. Use unused numeric IDs and an unused account name inside the container; map them to the intended host owner as appropriate.

For example, the published Alpine-based image can write files for host UID 1000 and GID 1000:

```bash
docker run --rm -v "$(pwd)/ssl:/ssl" \
  -e "CERT_DNS=lab.com" -e USER=certowner -e UID=1000 -e GID=1000 \
  soulteary/certs-maker:v3.8.2
```

This helper makes all output files readable by all users, including private keys. Restrict access to the output directory and restore restrictive key permissions when needed:

```bash
chmod 700 ssl
chmod 600 ssl/*.key
```

Apply those commands as the file owner. On native installations, use the operating system's file ownership tools rather than the container helper.
