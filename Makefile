satcors: $(shell find . -name "*.go")
	CC=$$(which musl-gcc) go build -ldflags='-s -w -linkmode external -extldflags "-static"' -o ./satcors

deploy: satcors
	ssh root@hulsmann 'systemctl stop satcors'
	scp satcors hulsmann:satcors/satcors
	ssh root@hulsmann 'systemctl start satcors'
