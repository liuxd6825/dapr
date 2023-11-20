go run /Users/lxd/go/src/github.com/liuxd6825/dapr/cmd/daprd/main.go  \
-log-level info \
-app-port 9037 \
-dapr-http-port 9038 \
-dapr-grpc-port 9039 \
-app-id duxm-scan-query-service \
--enable-metrics=false \
-config /Users/lxd/projects/duxm/duxm-scan/config-liuxd/dapr/config.yaml \
-components-path /Users/lxd/projects/duxm/duxm-scan/config-liuxd/dapr/components \
-placement-host-address=127.0.0.1:50005