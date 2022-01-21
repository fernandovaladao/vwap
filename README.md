
# Zero Hash Code Assessment: Volume-Weighted Average Price Calculation Engine
------------------------------------------------

This repository contains an implementation for a real-time VWAP (volume-weighted average price) calculation engine.

The project is a simple CLI tool that echos back the average price for the last 200 trading prices for a set of crypto currencies.

# Prerequisites

The only requirements to build and use this project are:
- `Docker` 18.09 or higher, since we're using `BuildKit` to build the containers' images; and
- `make`.

## macOS

* Install [Docker Desktop](https://www.docker.com/products/docker-desktop)
* Ensure that you have `make` (included with Xcode)

## Linux

* Install [Docker](https://docs.docker.com/engine/install/)
* Ensure that you have `make`

# Getting started

Building and installing the project can be accomplished using the following default `make` command:
```console
$ make
```
You will notice that this command:
* Built Docker images to run the application and to execute end-to-end tests; and
* Ran automated unit, integration and end-to-end tests to validate the solution.

```console
$ docker images
REPOSITORY                 TAG       IMAGE ID       CREATED        SIZE
zero-hash-vwap             v1.0.0    12cdcc046821   17 hours ago   6.39MB
zero-hash-vwap-e2e-tests   latest    41c718b2db0c   20 hours ago   721MB
```

You can then run the binary, as follows:
```console
$ make run
time="2022-01-21T11:38:25Z" level=error msg="Error message returned by trade stream client" message="Failed to subscribe" reason="ETH-USD is not a valid product"
time="2022-01-21T11:38:25Z" level=warning msg="Unknown trading message" trade="&{subscriptions    }"
time="2022-01-21T11:38:25Z" level=info trade_pair=BTC-USD vwap=194.12285
time="2022-01-21T11:38:25Z" level=info trade_pair=ETH-BTC vwap=0.00036705
time="2022-01-21T11:38:27Z" level=info trade_pair=BTC-USD vwap=387.48935
time="2022-01-21T11:38:31Z" level=info trade_pair=BTC-USD vwap=581.6122
time="2022-01-21T11:38:36Z" level=info trade_pair=BTC-USD vwap=774.3986
time="2022-01-21T11:38:41Z" level=info trade_pair=BTC-USD vwap=968.5214500000001
time="2022-01-21T11:38:46Z" level=info trade_pair=BTC-USD vwap=1161.9394
time="2022-01-21T11:38:52Z" level=info trade_pair=BTC-USD vwap=1356.06225
```
This command runs the `vwap` engine with this pre-defined set of trading pairs:
- BTC-USD,
- ETH-USD[^1], and
- ETH-BTC.

[^1]: Currently, the sandbox of `Coinbase Websocket` feed returns an error for `ETH-USD`, so no trading price is printed out for it. It is not clear if this is an intermitent issue with the service or a limitation.


You can inform an alternative set of trading pairs using the following variable:
```console
$ make run TRADE_PAIRS='BTC-USD BNB-USD'
```

Once you are done with running the project, you can type the following command to clean up the Docker images:
```console
$ make clean
Untagged: zero-hash-vwap-e2e-tests:latest
Deleted: sha256:41c718b2db0ca8bc1497c425c6a1c68a06e3add22046b06e84b749d460721a39
Untagged: zero-hash-vwap:v1.0.0
Deleted: sha256:12cdcc0468215b4db80f7d2098cf83a31c742583d4e67d53a45d9127c24aa742
```

# Structure of project

## Dockerfile

The [Dockerfile](./Dockerfile) codifies all the tools needed for the project
and the commands that need to be run for building, testing and running it.

## Makefile

The [Makefile](./Makefile) is purely used to script the required `docker build`
commands.

Besides those already mentioned on the top of this document, there are alternative targets to run specific suite of tests, such as:
```console
$ make unit-test
$ make integration-test
$ make e2e-test
```
If you would like to skip the tests, you can build and run the application using the following combination of targets:
```console
$ make build run
```

## Packages
### *data_structures*
This package contains a naive implementation of a circular queue using a fixed-size array, currently set as 200.

This data structure has *O(1)* space and time complexity for all of its operations, which made it to be very efficient in the scope of this problem.

### *storage_manager*
The *storage_manager*  package provides an interface to persist the last 200 trading prices for a
trading pair and also get the sum of these elements. It uses a circular queue as its buffer, which means
that all these operations run in *O(1)* time complexity and uses *O(1)* space.

The algorithm used to store new trading prices can be described as follows:
```python
StorageManager:
    buffer
    sum

func (StorageManager) store(price)
    if buffer.isFull():
        dequeue e from buffer
        decrement e in sum
    enqueue new price in buffer
    increment sum with new trading price
```

### *trade_streaming*
The purpose of this package is to provide an interface to stream trading prices for a set of trading pairs.
The current implementation feeds `Coinbase Websocket` in sandbox environment using `Gorilla WebSocket` api.

### *main*
This package wraps all the interfaces provided in the previous packages to build the `vwap` calculation engine. This means that this structure is composed of:
- a map between a trading pair and its storage manager;
- a client with a connection to the trading stream; and
- a log system, which currently prints out the calculated values to the standard console output.
The algorithm to calculate `vwap` for each trading pair is:

```python
VWAPEngine:
    storage_managers
    trade_client
    log
func (VWAPEngine) calculate:
    while true:
        trade = trade_client.conn.Read()
        sm = storage_managers[trade.pair]
        sm.store(trade.price)
        vwap = sm.sum() / 200
        log.Info(trade.pair, vwap)
```

The space complexity of this operation is *O(n)*, where *n* is the number of trading pairs being evaluated. The running complexity is hard
to guess, because we would have to rely on several variables we currently do not have information about, specially the time complexiy to read
value from the stream.