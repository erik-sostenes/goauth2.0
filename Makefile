include .env.example

server:
	bash ./scripts/build.bash

run: server
	./build