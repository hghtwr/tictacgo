build::
	GOOS=linux GOARCH=amd64 go build -o ./handler/connect ./handler/connect.go
	zip -j ./handler/connect.zip ./handler/connect

	GOOS=linux GOARCH=amd64 go build -o ./handler/turn ./handler/turn.go
	zip -j ./handler/turn.zip ./handler/turn

	GOOS=linux GOARCH=amd64 go build -o ./handler/disconnect ./handler/disconnect.go
	zip -j ./handler/disconnect.zip ./handler/disconnect
