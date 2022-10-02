prepare:
	mkdir -p /opt/htmls
	chmod 777 /opt/htmls

run: prepare
	go run serverWithChannels/main.go