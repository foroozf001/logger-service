# Logger Service

## Overview

The logger service provides an endpoint for storing logs in a mongodb back-end. It accepts HTTP REST on port 8080 and it accepts gRPC on tcp port 50051. In addition to this, it will expose custom metrics on an endpoint.

## Getting started

- Spin up the stack by executing `make up`
- Create artificial load using `hey` command-line tool

```bash
hey -z 5m -q 5 -m POST -H "Content-Type: application/json" -d '{"name": "example_log", "data": "this is an example log"}' http://localhost:8080/api/log
```

- View documents via mongo-express client at http://localhost:8081
- See custom prometheus metrics by calling the `/api/metrics` endpoint
  - `http_requests_total`
  - `response_status`
  - `http_response_time_seconds`

```bash
curl localhost:8080/api/metrics 
```
