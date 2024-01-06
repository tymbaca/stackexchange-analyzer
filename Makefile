run:
	@go run main.go

clickhouse:
	docker run -d --name some-clickhouse-server --ulimit nofile=262144:262144 -p 8123:8123 clickhouse/clickhouse-server
clickhouse-rm:
	docker stop some-clickhouse-server
	docker rm some-clickhouse-server
clickhouse-sh:
	docker exec -it some-clickhouse-server clickhouse-client

grafana:
	docker run -d -p 3000:3000 --name grafana grafana/grafana
grafana-rm:
	docker stop grafana
	docker rm grafana
