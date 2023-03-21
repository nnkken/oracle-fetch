#!/usr/bin/make -f

UID := $(shell id -u)
GID := $(shell id -g)
PG_TEST_PORT ?= 5433
CHAINLINK_ETH_PATH := datasource/chainlink-eth/contract

DOCKER_RUN := docker run --rm -u ${UID}:${GID}

gen-chainlink-eth: ${CHAINLINK_ETH_PATH}/contract.go

${CHAINLINK_ETH_PATH}/contract.go: ${CHAINLINK_ETH_PATH}/AggregatorV3Interface.abi
	${DOCKER_RUN} -v $(PWD)/${CHAINLINK_ETH_PATH}:/contract -w /contract ethereum/client-go:alltools-release-1.11 abigen --abi /contract/AggregatorV3Interface.abi --pkg contract --out /contract/contract.go

${CHAINLINK_ETH_PATH}/AggregatorV3Interface.abi: ${CHAINLINK_ETH_PATH}/AggregatorV3Interface.sol
	${DOCKER_RUN} -v $(PWD)/${CHAINLINK_ETH_PATH}:/contract -w /contract ethereum/solc:0.8.19 --abi --overwrite -o /contract /contract/AggregatorV3Interface.sol

gen-mock:
	${DOCKER_RUN} -v $(PWD):/src -w /src vektra/mockery:v2 --name "ChainLinkContract|DataSource"

setup-test-db:
	docker run --rm -p "${PG_TEST_PORT}:5432" -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres_test postgres:15

# TODO: containerize this, currently needs `go install github.com/swaggo/swag/cmd/swag@latest`
gen-swag:
	swag init

test:
	PG_TEST_DB_URL="postgres://postgres:postgres@localhost:${PG_TEST_PORT}/postgres_test" go test -count 1 -p 1 -v ./...

.PHONY: gen-chainlink-eth gen-mock test-db test
