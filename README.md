# Hetty

Hetty is an HTTP toolkit for security research. It aims to become an open source
alternative to commercial software like Burp Suite Pro, with powerful features
tailored to the needs of the infosec and bug bounty community.

![Screenshot](docs/screenshot.png)

## Features

- Man-in-the-middle (MITM) HTTP/1.1 proxy, with logs
- Project management (e.g. for grouping logs and scope)
- Search and filter requests/responses
- Intercept requests and responses, and modify before forwarding
- Sender module for manually crafting/repeating HTTP requests

## Requirements

Hetty is a single binary that bundles the API server (Go) and web interface (Next.js).

To use Hetty, you need to:

1. Have a local CA certificate trusted by your browser/OS
2. Configure your browser or system to use Hetty as an HTTP proxy

## Installation

### Pre-built binaries

Pre-built binaries are available on the [releases](https://github.com/dstotijn/hetty/releases) page.

### Docker

```sh
docker run -v $HOME/.hetty:/root/.hetty -p 8080:8080 dstotijn/hetty
```

### Build from source

```sh
git clone https://github.com/dstotijn/hetty.git
cd hetty
make build
```

## Usage

```
Usage of hetty:
  -addr string
        TCP address to listen on, in the form "host:port" (default ":8080")
  -adminPath string
        File path to admin build directory (default uses embedded assets)
  -cert string
        File path to root CA certificate (default "~/.hetty/hetty_cert.pem")
  -key string
        File path to root CA private key (default "~/.hetty/hetty_key.pem")
  -db string
        File path to database (default "~/.hetty/hetty.db")
  -version
        Print version and exit
```

### CA Certificate Setup

On first run, Hetty will generate a root CA certificate and private key if they
don't exist. You need to import the certificate into your browser/OS trust store.

The certificate is stored at `~/.hetty/hetty_cert.pem` by default.

> **Note (personal):** On macOS, you can trust the cert quickly with:
> `sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.hetty/hetty_cert.pem`

### Proxy Configuration

Configure your browser or system to use `http://localhost:8080` as the HTTP proxy.

> **Note (personal):** In Firefox, go to Settings → General → Network Settings → Manual proxy
> configuration. Set HTTP Proxy to `127.0.0.1` and Port to `8080`. Make sure to check
> "Also use this proxy for HTTPS". Firefox has its own certificate store, so you'll
> also need to import the CA cert under Settings → Privacy & Security → Certificates → View Certificates.

> **Note (personal):** On Linux with Chrome/Chromium, you can launch the browser already pointed at
> the proxy with:
> `chromium --proxy-server="http://127.0.0.1:8080" --ignore-certificate-errors-spki-list`
> Though it's cleaner to import the CA cert and skip `--ignore-certificate-errors-spki-list`.

## Development

### Prerequisites

- Go 1.21+
- Node.js 18+
- pnpm

### Running locally

```sh
# Start the Next.js dev server
cd admin
pnpm install
pnpm dev

# In another terminal, start the Go server
go run ./cmd/hetty -adminPath ./admin/dist
```

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for
guidelines on how to contribute to this project.

## License

[Apa