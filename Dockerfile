FROM golang:1.17-alpine3.15 as build-env

# 环境
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY="https://goproxy.io,direct"

# 事先复制过来可以利用docker的缓存 只要mod文件不动就可以不重新下载依赖
COPY go.mod /app/
COPY go.sum /app/
WORKDIR /app
RUN go mod download

COPY . /app
RUN  go build -a -installsuffix cgo -o run main.go


FROM alpine:3.15

# 换源
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.12/main/ > /etc/apk/repositories
# 安装mysql-client, 时区
RUN apk add --no-cache mysql-client tzdata
ENV TZ Asia/Shanghai

COPY --from=build-env /app/run /

CMD ./run

# docker build -t daihaoxiaofei/mysqldump:1 .
# docker push daihaoxiaofei/mysqldump:1

# docker tag 34d2ab2e2dcc mysqldump:v1 # 改标签
# docker rmi backups:v1 # 删除