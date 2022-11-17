# 一个备份mysql数据库的镜像,可备份多个数据库

### 使用:

    按照 config.yaml.example 格式复制一份改名为:config.yaml 放到本目录使用
    建议放到docker中使用 已打包好镜像:daihaoxiaofei/mysqldump:1
    也可自行编译后放到有mysql-client的环境下直接运行

### docker-compose.yml:

    version: '3'
    services:
        backup:
            container_name: backup
            image: daihaoxiaofei/mysqldump:1
            restart: always
            volumes:
                - ./backup:/backup # 备份文件
                - ./config/config.yaml:/config/config.yaml  # 配置文件

### 原理:

    在alpine镜像中安装了mysqldump(mysql-client)用于备份
    可以用yaml文件配置多个需要备份的数据库
    用go写了备份文件的规则:
        频率:  yaml中自定义配置 默认= [0,59),* * * *
        清理:  可保留时间点:当前小时,前2 4 8 16小时,
                当日 前一天 前2天凌晨6点
                当月1号凌晨6点 前一月 前2月 前3月 凌晨6点
                其余会每小时清理
        注: 会根据自身要求随时更改 如有需要请自行前往: github 修改源码重新编译



