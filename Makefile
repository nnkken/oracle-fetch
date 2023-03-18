#!/usr/bin/make -f

UID := $(shell id -u)
GID := $(shell id -g)
CHAINLINK_ETH_PATH := datasource/chainlink-eth/contract

gen-chainlink-eth: ${CHAINLINK_ETH_PATH}/contract.go

${CHAINLINK_ETH_PATH}/contract.go: ${CHAINLINK_ETH_PATH}/AggregatorV3Interface.abi
	docker run --rm -v $(PWD)/${CHAINLINK_ETH_PATH}:/contract -w /contract -u ${UID}:${GID} ethereum/client-go:alltools-release-1.11 abigen --abi /contract/AggregatorV3Interface.abi --pkg contract --out /contract/contract.go

${CHAINLINK_ETH_PATH}/AggregatorV3Interface.abi: ${CHAINLINK_ETH_PATH}/AggregatorV3Interface.sol
	docker run --rm -v $(PWD)/${CHAINLINK_ETH_PATH}:/contract -w /contract -u ${UID}:${GID} ethereum/solc:0.8.19 --abi --overwrite -o /contract /contract/AggregatorV3Interface.sol

.PHONY: gen-chainlink-eth
