go run /Users/lxd/go/src/github.com/liuxd6825/dapr/cmd/daprd/main.go \
-log-level debug  \
-app-port 35010 \
-dapr-http-port 35011 \
-dapr-grpc-port 35012 \
-app-id duxm-fund-record-service-command-service \
--enable-metrics=false \
-config /Users/lxd/projects/duxm/duxm-fund-record-service/config/dapr/config.yaml \
-components-path /Users/lxd/projects/duxm/duxm-fund-record-service/config/dapr/components \
-placement-host-address=127.0.0.1:50005
