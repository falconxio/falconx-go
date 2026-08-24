module github.com/falconxio/falconx-go

go 1.14

// TODO: change this to falconx repo once the fork is done successfully
replace github.com/graarh/golang-socketio => github.com/pradeepfx/golang-socketio v0.0.0-20230424110355-180e43d5e4f1

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/graarh/golang-socketio v0.0.0-20170510162725-2c44953b9b5f
	go.uber.org/multierr v1.7.0
)
