# Nebula-API

[![CI](https://github.com/synchthia/nebula-api/workflows/CI/badge.svg?branch=master&event=push)](https://github.com/synchthia/nebula-api/actions?query=workflow%3ACI)

STARTAIL network management system.

## Environment Variables

| Environment Variables     | Description           | Default                                                                           |
| ------------------------- | --------------------- | --------------------------------------------------------------------------------- |
| `MYSQL_CONNECTION_STRING` | MySQL address         | `root:docker@tcp(localhost:3306)/nebula?charset=utf8mb4&parseTime=True&loc=Local` |
| `REDIS_ADDRESS`           | Redis address         | `localhost:6379`                                                                  |
| `GRPC_LISTEN_PORT`        | gRPC Listening port   | `:17200`                                                                          |
| `ENABLE_IP_FILTER`        | db-ip.com IP checker  | false                                                                             |
| `DB_IP_TOKEN`             | db-ip.com Private Key | none                                                                              |
| `DEBUG`                   | Enable debug output   | none                                                                              |
