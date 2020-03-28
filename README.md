# Nebula-API

[![CI](https://github.com/synchthia/nebula-api/workflows/CI/badge.svg?branch=master&event=push)](https://github.com/synchthia/nebula-api/actions?query=workflow%3ACI)

STARTAIL network management system.

## Environment Variables

| Environment Variables     | Description                         | Default                     |
| ------------------------- | ----------------------------------- | --------------------------- |
| `MONGO_CONNECTION_STRING` | MongoDB address (connection string) | `mongodb://localhost:27017` |
| `REDIS_ADDRESS`           | Redis address                       | `localhost:6379`            |
| `GRPC_LISTEN_PORT`        | gRPC Listening port                 | `:17200`                    |
| `DEBUG`                   | Enable debug output                 | none                        |
