FROM golang:1.19.8-alpine3.17 as build-env

# 启用go module
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPATH=/GOPATH \
    GOPROXY="https://goproxy.cn,direct"

# 事先复制过来可以利用docker的缓存 只要mod文件不动就可以不重新下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 编译
COPY . .
RUN go build -a -installsuffix cgo -o run main.go


FROM alpine:3.17

# 换源 安装mysql-client, 时区
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.12/main/ > /etc/apk/repositories  &&\
    apk add --no-cache mysql-client tzdata &&\
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=build-env /go/run /

CMD ./run

# docker build -t daihaoxiaofei/mysqldump:1 .
# docker push daihaoxiaofei/mysqldump:1

# docker tag 34d2ab2e2dcc mysqldump:v1 # 改标签
# docker rmi backups:v1 # 删除
