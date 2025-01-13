# DAPR编译
## mac 准备
下载 Command Line Tools
https://developer.apple.com/download/all/?q=Command%20Line%20Tools

	$ rm -rf /Library/Developer/CommandLineTools 
	$ xcode-select --install

### 安装docker私仓
假设：私仓ip为192.168.64.12
$ docker run -d  -p 5000:5000  --restart=always --name registry  -v /opt/data/registry:/var/lib/registry registry

### 安装Helm
1. 下载[所需版本]
   https://github.com/helm/helm/releases
2. 打开包装
```shell
$ tar -zxvf helm-v3.0.0-linux-amd64.tgz
```
3. helm在解压后的目录中找到二进制文件，然后将其移至所需的目标位置
```shell
   $ mv linux-amd64/helm /usr/local/bin/helm
```

4. 在客户端运行：
```shell
    $ helm help
``` 
5. 配置国内Chart仓库
   微软仓库（http://mirror.azure.cn/kubernetes/charts/）这个仓库强烈推荐，基本上官网有的chart这里都有。
   阿里云仓库（https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts ）
   官方仓库（https://hub.kubeapps.com/charts/incubator）官方chart仓库，国内有点不好使

6. 添加存储库
```shell
helm repo add stable http://mirror.azure.cn/kubernetes/charts
helm repo add aliyun https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts
helm repo update
``` 
### 3.创建命名空间
```shell
$ kubectl create namesapce dapr-system
``` 
### 4.设置环境变量
```shell
$ export DAPR_REGISTRY=192.168.64.12:5000
$ export DAPR_TAG=dev
$ export  TARGET_OS=linux
$ export  TARGET_ARCH=arm64    //arm架构cpu
``` 
### 5.编译Docker镜像文件
进入dapr项目目录
```shell
$ cd dapr
``` 

### 编译二进制文件
设置版本号：REL_VERSION=v1.15-250111
```shell
$ make build GOOS=linux    GOARCH=arm64  REL_VERSION=v1.15-250111
$ make build GOOS=windows  GOARCH=amd64  REL_VERSION=v1.15-250111
$ make build GOOS=darwin   GOARCH=arm64  REL_VERSION=v1.15-250111
```

修改docker/docker.rm文件， 增加参数 --load
```shell
docker-build: check-docker-env check-arch 
    $(info Building $(DOCKER_IMAGE_TAG) docker image ...) 
ifeq ($(TARGET_ARCH),amd64) 
    $(DOCKER) build --load  PKG_FILES=* -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
    $(DOCKER) build --load  PKG_FILES=daprd -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_RUNTIME_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
    $(DOCKER) build --load  PKG_FILES=placement -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_PLACEMENT_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
    $(DOCKER) build --load  PKG_FILES=sentry -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_SENTRY_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
else
    -$(DOCKER) buildx create --use --name daprbuild
    -$(DOCKER) run --rm --privileged multiarch/qemu-user-static --reset -p yes
    $(DOCKER) buildx build --load  --build-arg PKG_FILES=*         --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
    $(DOCKER) buildx build --load  --build-arg PKG_FILES=daprd     --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_RUNTIME_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
    $(DOCKER) buildx build --load  --build-arg PKG_FILES=placement --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_PLACEMENT_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
    $(DOCKER) buildx build --load  --build-arg PKG_FILES=sentry    --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_SENTRY_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
endif
```



### 生成docker image
```shell
$ sudo make docker-build TARGET_OS=linux TARGET_ARCH=arm64 DAPR_REGISTRY=192.168.64.12 DAPR_TAG=dapr
```

### 查看是否生成
```shell
$ docker images

REPOSITORY                              TAG                                   IMAGE ID              CREATED                SIZE
192.168.64.12/sentry                 dapr-linux-arm64             8f909cf60e46        13 minutes ago      37MB
192.168.64.12/placement          dapr-linux-arm64             04eb59389523      13 minutes ago      16.3MB
192.168.64.12/daprd                 dapr-linux-arm64             936608e34a01       13 minutes ago      105MB
192.168.64.12/dapr                   dapr-linux-arm64             207d8890f756       13 minutes ago      286MB

```


### 推送image文件到私库中
```shell
$ docker push 192.168.64.12/sentry:dapr-linux-arm64
$ docker push 192.168.64.12/placement:dapr-linux-arm64
$ docker push 192.168.64.12/daprd:dapr-linux-arm64
$ docker push 192.168.64.12/dapr:dapr-linux-arm64
```


/*
推送image到私仓
$ sudo make docker-push TARGET_OS=linux TARGET_ARCH=arm64 DAPR_REGISTRY=192.168.64.12  DAPR_TAG=dapr
*/

/*
标记docker
$ docker tag dapr:dev-linux-arm64 192.168.64.12:5000/dapr:dev-linux-arm64
或
$ docker tag dapr:dev-linux-arm64 192.168.64.12:5000/dapr:dev-linux-amd64

推到push中
$ docker push 192.168.64.12:5000/dapr:dev-linux-arm64
或
$ docker push 192.168.64.12:5000/dapr:dev-linux-amd64
*/

### 查看私仓
```shell
$ curl 192.168.64.12:443/v2/_catalog  
$ curl 192.168.64.12:443/v2/dapr/tags/list
```



## 执行make部署

### 进入项目目录
```shell
$ cd "dapr项目根目录"
```

### 删除原有安装
```shell
$ dapr uninstall -k
```

### 安装方法
```shell
$ make docker-deploy-k8s TARGET_OS=linux TARGET_ARCH=arm64 DAPR_REGISTRY=192.168.64.12  DAPR_TAG=dapr
```

## 查询pod
```shell
$ kubectl get pod -n dapr-system
NAMESPACE  NAME                                                         READY             STATUS              RESTARTS         AGE
dapr-system   dapr-dashboard-c8dd8d969-xgdtx             1/1                 Running                   0                  18m
dapr-system   dapr-operator-7c7d5f887b-b5dsn              0/1                 ImagePullBackOff    0                  18m
dapr-system   dapr-placement-server-0                             0/1                 ImagePullBackOff    0                  18m
dapr-system   dapr-sentry-6cc6687cc8-sqrcr                     0/1                 ImagePullBackOff    0                  18m
dapr-system   dapr-sidecar-injector-698756d5d7-46fpd   0/1                 ImagePullBackOff    0                  18m
```

## 解决K8s的 ImagePullBackOff
这个问题说明下载image时出错，image文件没有找到。检查docker私仓镜像名称与pull时是否一致。
$ kubectl get pods dapr-operator-7c7d5f887b-b5dsn -n dapr-system -o yaml | grep image:

## 解决K8s的 Pending
$ kubectl get pod -n dapr-system

查看具体原因
$ kubectl describe pod <pod> -n dapr-system

可能性一：生成平台与部署平台不一致。

	
	