[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://godoc.org/github.com/falconxio/falconx-go)

FalconX Go Client
==================================
Go Client to connect to FalconX Rest/WebSocket APIs. Library contains clients for rest and websocket interfaces, and example files for using them.

[Comprehensive Rest API Docs](https://falconx.io/docs#rest-api)

[Comprehensive Socket IO API Docs](https://falconx.io/docs#websocket-api)


Installing and building the library
===================================

This project requires Go 1.14. As of the time of writing.
To use this package in your own code, make sure your `GO_PATH` environment variable is correctly set, and install it using `go get`:

    go get github.com/falconxio/falconx-go

Then, you can include it in your project( Please refer to [client_examples](https://github.com/falconxio/falconx-go/tree/main/client_examples) for detailed examples of usage.)

	import "github.com/falconxio/falconx-go/clients"

Alternatively, you can clone it yourself:

    git clone github.com/falconxio/falconx-go.git


How to run examples?
==================================
```
go run run_examples -api_key=XXX -secret=XXX -passphrase=XXX -example_set=rest
go run run_examples -api_key=XXX -secret=XXX -passphrase=XXX -example_set=websocket
```

Questions?
==================================
In case of any questions please contact support@falconx.io



