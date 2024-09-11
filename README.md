
# Go REST HTTP Server with Rate Throttling over GRPC

A Simple REST API web-server tracking IP-Addresses and request-count with rate throttling over GRPC via a client, exposed metrics via Prometheus and visualisation through Grafana




## Development

Make sure Redis is up and running on localhost:6379

```bash
  go run main.go
```

Run prometheus via prometheus.yml

```bash
docker run \
    -p 9090:9090 \
    -v /path/to/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus
```

Choose Prometheus as a Data Source in Grafana to start visualising.