Redis Lite

A simple Redis-like TCP server and client in Go.
Features

    Key-Value Store: set and get with optional expiration.
    Concurrent Clients: Handles multiple connections.

Usage

    Setup:

    bash

git clone https://github.com/abdoroot/redis-lite.git
cd redis-lite
make build

Run Server:

bash

make run

Commands:

    Set: make set key=mykey value=myvalue [expire=10]
    Get: make get key=mykey