migrations:
	bash ./scripts/migrations.bash

keys: 
	bash ./scripts/keys.bash

server:
	bash ./scripts/build.bash

run: server
	./build