Redis Lite - A Simple Redis-like TCP Server in Go

Redis Lite is a lightweight, Redis-inspired TCP server and client written in Go. It supports basic set and get operations with optional expiration times, using a simple custom protocol.
Features

    Basic Commands: Implements set, get with optional expiration times.
    Concurrent Client Handling: Efficiently manages multiple client connections.
    Custom Protocol: Uses a straightforward text-based protocol for communication.
    Timeout Handling: Read and write deadlines to prevent indefinite blocking.

Getting Started
Prerequisites

    Go 1.16 or later
    Make

Installation

    Clone the Repository:

    bash

git clone https://github.com/abdoroot/redis-lite.git
cd redis-lite

Build the Server and Client:

bash

    make build

Usage

    Start the Server:

    bash

make run

This will start the server on the default port 8080.

Set a Key-Value Pair with Optional Expiration:

bash

make set key=mykey value=myvalue [expire=5]

The optional expire argument sets the key's expiration time in seconds.

Get the Value for a Key:

bash

    make get key=mykey

Project Structure

    server/: Contains the TCP server code.
    client/: Contains the TCP client code.
    protocol.go: Defines the custom protocol for command communication.

Example

    Start the Server:

    bash

make run

Output:

csharp

tcp server is running on port :8080

Set a Key-Value Pair:

bash

make set key=mykey value=myvalue

Output:

diff

+OK

Get the Value for a Key:

bash

make get key=mykey

Output:

myvalue

Set a Key-Value Pair with Expiration:

bash

    make set key=temporary value=expiring expire=10

    This will set the key temporary to expire in 10 seconds.

Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.
License

This project is licensed under the MIT License - see the LICENSE file for details.
