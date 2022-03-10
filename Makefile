.PHONY: build dry systemd install

all: dry

BINARY_NAME=netflix-agent
BINARY_PATH=/usr/local/bin/

build:
	go build -o ${BINARY_NAME} main.go

dry:
	go build -o ${BINARY_NAME} main.go

systemd: build
	cp ./{BINARY_NAME} ${BINARY_PATH}
	cp ./netflix-agent.service /lib/systemd/system/
	if [ ! -d /etc/netflix-agent ]; then
		mkdir /etc/netflix-agent
	fi
	cp ./config.yaml /etc/netflix-agent/

install: systemd
	systemctl daemon-reload
	systemctl enable netflix-agent
	systemctl start netflix-agent
