go run /Users/lxd/go/src/github.com/liuxd6825/dapr/cmd/daprd/main.go \
-app-port 9010 \
-dapr-http-port 9011 \
-dapr-grpc-port 9012 \
-app-id duxm-graph-command-service \
-enable-metrics=false \
-config /Users/lxd/projects/duxm/duxm-graph/config/dapr/config.yaml \
-components-path /Users/lxd/projects/duxm/duxm-graph/config/dapr/components \
-placement-host-address 127.0.0.1:50005