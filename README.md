# SPOT

Single Point of Testing

## TODO

See [project](https://github.com/users/vmarlier/projects/1)

## API Design Overview

This service follows a REST-RPC hybrid approach with no data storage. Each request is processed immediately, and responses are returned without persistence.

Endpoints:
- POST /v1/tests/execute – Runs a network test with the given parameters.
- POST /v1/requests/trigger – Sends a single request and returns the response.
- POST /v1/requests/batch/trigger – Sends multiple requests based on defined parameters.

All operations are stateless. If logging or result history is needed, it must be handled externally.
