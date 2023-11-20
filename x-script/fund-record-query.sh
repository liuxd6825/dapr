go run /Users/lxd/go/src/github.com/liuxd6825/dapr/cmd/daprd/main.go  \
-log-level debug  \
-app-port 35020 \
-dapr-http-port 35021 \
-dapr-grpc-port 35022 \
-app-id duxm-fund-record-service-query-service \
--enable-metrics=false \
-config /Users/lxd/projects/duxm/duxm-fund-record-service/config/dapr/config.yaml \
-components-path /Users/lxd/projects/duxm/duxm-fund-record-service/config/dapr/components