# Oracle Fetch

This is a simple tool to fetch price data from oracles like Chainlink.

It runs periodically and fetches the latest price from the oracle and stores it into database.

It also provides RESTful API to query the price data.

## Supported Oracles

Currently only Chainlink on Ethereum is supported.

## Usage

The tool is split into 2 sub-commands: `fetch` and `api`.

Common parameters:

 - `database-url`: the Postgres database URL for connecting the database. Format: `postgres://USERNAME:PASSWORD@HOST:PORT/DB_NAME`
 - `log-level`: the log level: `debug | info | warn | error | dpanic | panic | fatal`. Default: `info`.
 - `log-format`: the log format: `json | console`. Default: `console`.
 - `log-output`: the log output: `stdout | stderr | /somewhere/to/some/file`. Can be set multiple times for multiple outputs. Default: `stderr`.

### fetch

`oracle-fetch fetch --database-url <postgres-database-url> --config <config.json>`

 - `config`: the configuration file. See `config.json` as example, which is for Chainlink Ethereum mainnet, with rate limiter setup.

### api

`oracle-fetch api --database-url <postgres-database-url> --host <host-with-port>`

 - `host`: the host and port to listen on. `:8080` means listen on port 8080 on all interfaces, while `localhost:8080` means listen on port 8080 on only localhost connections.

## APIs

See `docs/swagger.yaml` or `docs/swagger.json` for details. You may also access `/swagger/index.html` for documents with Swagger UI.

## Development

You may use Docker Compose to easily start a development environment, with Postgres database setup.

In `docker-compose.yml`, configurate the `--eth-endpoint` part under `fetcher` service, then run `docker compose up` to start the service, and `docker compose build` to rebuild.

## Testing

To run test cases, you need an empty Postgres database. The fixture library only allows connecting to database with `test` suffix in the name for safety.

You may run `make setup-test-db` to setup a test database, and `make test` to run test cases.

Note that by default the testing and the test database are connected through port `5433` to distinguish with usual production port (`5432`). You may change the `PG_TEST_PORT` environment variable for a different port.

## Deployment

For deployment, you should set `GIN_MODE` environment variable to `release` to disable debug mode on API gateway.
