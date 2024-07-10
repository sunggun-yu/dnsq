# dnsq (DNS Query)

`dnsq` is a simple DNS lookup tool written in Go. It provides both a command-line interface and a REST API server for performing simple CNAME, A, and AAAA DNS queries with given multiple host names.

- [Installation](#installation)
  - [HomeBrew](#homebrew)
  - [Docker](#docker)
- [Usage](#usage)
  - [CLI](#cli)
  - [Server](#server)
    - [API Endpoint](#api-endpoint)
  - [Kubernetes Deployment](#kubernetes-deployment)

## Installation

### HomeBrew

```bash
brew install sunggun-yu/tap/dnsq
```

### Docker

```bash
docker pull ghcr.io/sunggun-yu/dnsq
```

## Usage

### CLI

```bash
dnsq google.com www.github.com aws.amazon.com www.facebook.com
```

This will return DNS information for multiple hosts google.com, www.github.com, aws.amazon.com, and www.facebook.com

### Server

```bash
dnsq server
```

By default, the server runs on port 8080. You can specify a different port using the `--port` or `-p` flag.

```bash
dnsq server -p 9090
```

you can also use container image for server instance.

```bash
docker run -d --name dnsq-server \
  -p 8080:8080 \
  ghcr.io/sunggun-yu/dnsq-server
```

#### API Endpoint

The API exposes a single endpoint for DNS lookups:

- Endpoint: `/api/lookup`
- Method: `GET`
- Query Parameter: `hosts` (comma-separated list of domain names)

Example request:

```bash
curl -s http://localhost:8080/api/lookup?hosts=google.com,blog.meowhq.dev,www.facebook.com
```

Example response: (assuming piped with `jq` for pretty json)

```json
{
  "blog.meowhq.dev": [
    {
      "host": "blog.meowhq.dev",
      "type": "A",
      "data": "76.76.21.21"
    }
  ],
  "google.com": [
    {
      "host": "google.com",
      "type": "A",
      "data": "142.250.80.110"
    }
  ],
  "www.facebook.com": [
    {
      "host": "www.facebook.com",
      "type": "A",
      "data": "31.13.71.36"
    }
  ]
}
```

### Kubernetes Deployment

```bash
kubectl create ns dnsq-server
kubectl apply -f https://raw.githubusercontent.com/sunggun-yu/dnsq/main/manifests/install.yaml
```
