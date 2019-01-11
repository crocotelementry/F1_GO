# F1_GO

> A Live Telemetry Dashboard, Storage, and Analyzer for Codemasters F1 2018 game for PC, XBOX, and Playstation

F1_GO is Written in Go and Utilizes Websockets, Redis, and MYSQL

---------------------------------------
  * [Features](#features)
  * [Requirements](#requirements)
  * [Installation](#installation)
  * [Usage](#usage)
  * [License](#license)

---------------------------------------

## Features
  * None
  * So
  * Far
  * Lol
  * Hope
  * This
  * Works?

## Requirements
  * Some version of go so you know, stuff actually works
  * Gorrila Web Toolkit Websocket
  * Gorilla Web Toolkit Mux
  * Redigo, a Go client for the Redis database.
  * Go-MySQL-Driver, A MySQL-Driver for Go's database/sql package

---------------------------------------

## Installation
Simple install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go tool](https://golang.org/cmd/go/ "go command") from shell:
```bash
$ go get -u github.com/crocotelementry/F1_GO
```
Make sure [Git is installed](https://git-scm.com/downloads) on your machine and in your system's `PATH`.

Until we find a way to have our requirements including in the F1_GO package, We will also need to install four more items into your go path.

**gorilla/websocket:**
```bash
$ go get github.com/gorilla/websocket
```

**gorilla/mux:**
```bash
$ go get github.com/gorilla/mux
```

**redigo:**
```bash
$ go get github.com/gomodule/redigo/redis
```

**go-sql-driver:**
```bash
$ go get -u github.com/go-sql-driver/mysql
```

## Usage
*F1_GO* is ran by running the main executable. Some features that are critical to *F1_GO's* usability are able to be ran from the terminal window in which you
start *F1_GO*, but it is not recommended. After *F1_GO* is started, all that is needed is to access the websocket from a web browser at the following address: *http://localhost:8080/*

To run *F1_GO*:
```go
go run *.go
```

---------------------------------------

## License
F1_GO is licensed under the [MIT License](https://raw.github.com/crocotelementry/F1_GO/master/LICENSE)

MIT License summarizes the license scope as follows:
> A short and simple permissive license with conditions only requiring preservation of copyright and license notices. Licensed works, modifications, and larger works may be distributed under different terms and without source code.


That means:
  * Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:

  * The above copyright notice and this permission notice shall be included in all
  copies or substantial portions of the Software.

  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  SOFTWARE.

You can read the full terms here: [LICENSE](https://raw.github.com/crocotelementry/F1_GO/master/LICENSE).
