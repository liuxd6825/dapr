go run /Users/lxd/go/src/github.com/liuxd6825/dapr/cmd/daprd/main.go \
--log-level debug  \
-app-port 9017 \
-dapr-http-port 9018 \
-dapr-grpc-port 9019 \
-app-id duxm-graph-query-service \
-enable-metrics=false \
-config /Users/lxd/projects/duxm/duxm-graph/config/dapr/config.yaml \
-components-path /Users/lxd/projects/duxm/duxm-graph/config/dapr/components \
-placement-host-address=127.0.0.1:-port=50005