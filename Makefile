build::
	GOOS=linux GOARCH=amd64 go build -o ./handler/connect ./handler/connect/connect.go
	zip -j ./handler/connect.zip ./handler/connect/connect
	rm ./handler/connect/connect

	GOOS=linux GOARCH=amd64 go build -o ./handler/turn ./handler/turn/turn.go
	zip -j ./handler/turn.zip ./handler/turn/turn
	rm ./handler/turn/turn

	GOOS=linux GOARCH=amd64 go build -o ./handler/disconnect ./handler/disconnect/disconnect.go
	zip -j ./handler/disconnect.zip ./handler/disconnect/disconnect
	rm ./handler/disconnect/disconnect

	GOOS=linux GOARCH=amd64 go build -o ./handler/setup ./handler/setup/setup.go
	zip -j ./handler/setup.zip ./handler/setup/setup
	rm ./handler/setup/setup