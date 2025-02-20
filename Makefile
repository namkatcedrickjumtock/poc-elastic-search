run:
	./elastic-start-local/start.sh
	go run ./cmd/main.go

stop:
	./elastic-start-local/stop.sh

install-tools:
	go get github.com/elastic/go-elasticsearch/v8@latest

uninstall:
	./elastic-start-local/uninstall.sh

# init commands pulls elastic & kibana images.
init:
	curl -fsSL https://elastic.co/start-local | sh

connector:
	docker run \
	-v "$${HOME}/elastic-connectors:/config" \
	--tty \
	--rm \
	docker.elastic.co/integrations/elastic-connectors:8.17.1 \
	/app/bin/elastic-ingest \
	-c /config/config.yml