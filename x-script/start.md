
## hserver-dapr
```shell
daprd -log-level=info
    -app-port=43001
    -dapr-http-port=43002
    -dapr-grpc-port=43003
    -app-id=duxm-master-cmd-service
    -config=/Users/lxd/Projects/duxm/h-master/server/config/dapr/dev_lxd/config.yaml
    -resources-path=/Users/lxd/Projects/duxm/h-master/server/config/dapr/dev_lxd/cmd-components
    -metrics-port=9191
    -placement-host-address=127.0.0.1:50005
    -scheduler-host-address=127.0.0.1:50006
```
