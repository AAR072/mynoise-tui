APP_NAME = mynoise-tui
SRC_DIR = ./src

.PHONY: build run clean
.DEFAULT_GOAL := run

build:
	mkdir -p bin
	cd src && go build -o ../bin/mynoise-tui

run:
	cd src && go run main.go

clean:
	rm -rf bin
