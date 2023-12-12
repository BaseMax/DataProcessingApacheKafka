KAFKA_SERVER=localhost:29092

USER := admin
PASSWORD := password
GRAFANA_SERVER := localhost:3000
DSDIR := grafana/datasources

all: run

run: down
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose build

grafana_import_ds:
	for i in ${DSDIR}/*; do \
		curl -X "POST" "http://${GRAFANA_SERVER}/api/datasources" \
			-H "Content-Type: application/json" \
			--user ${USER}:${PASSWORD} \
			--data-binary @$$i; \
	done

grafana_export_ds:
	curl -s "http://${GRAFANA_SERVER}/api/datasources" -u ${USER}:${PASSWORD} | jq -c -M '.[]' | split -l 1 - ${DSDIR}/

cleanup:
	rm -rf ${SAMPLE_OUT} *_data
