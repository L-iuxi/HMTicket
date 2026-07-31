# Docker
### Docker服务指令
```bash
#启动docker服务
systemctl start docker

#关闭docker服务
systemctl stop docker

#查看docker服务状态
systemctl status docker

#设置开机启动
systemctl enable docker
```

### Docker镜像指令
```bash
#查看本地镜像
docker images    

#拉取远程镜像
docker pull name:version

#搜索远程镜像
docker search name

#删除本地镜像
docker rmi imageid
```


### Docker容器指令
```bash
#创建一个新容器并启动
#参数分别为 持续运行|独立终端|名字|镜像及版本号|配置数据卷|端口映射|进入容器的初始化指令
docker run -it --name= repository:version  -v -p port:port cmd

#查看正在运行的容器
docker ps

#查看历史容器
docker ps -a

#进入容器
docker exec -it name cmd

#启动已经存在的容器
docker start dockerid

#停止容器
docker stop name

#删除容器
#正在运行的容器不能删除
docker rm name

#查看容器信息
docker inspect name

#容器转镜像
docker commit dockerid imagename version
```

### Dockerfile 关键字总览

| 关键字 | 作用 | 执行阶段 | 是否生成镜像层 | 常见使用场景 |
|---|---|---|---|---|
| FROM | 指定基础镜像 | 构建阶段 | 是 | 指定 Linux、Go、Java、Node 等运行环境 |
| RUN | 执行 Linux 命令 | 构建阶段 | 是 | 安装软件、编译代码、修改环境 |
| COPY | 复制文件到镜像 | 构建阶段 | 是 | 复制代码、配置文件、资源文件 |
| ADD | 复制文件（增强版 COPY） | 构建阶段 | 是 | 自动解压 tar 包、复制远程文件 |
| WORKDIR | 设置工作目录 | 构建阶段/运行阶段 | 否 | 后续 RUN、COPY、CMD 的默认目录 |
| ENV | 设置环境变量 | 构建阶段/运行阶段 | 否 | 配置数据库地址、环境参数 |
| ARG | 定义构建参数 | 构建阶段 | 否 | 构建镜像时传入变量 |
| EXPOSE | 声明容器监听端口 | 构建阶段 | 否 | 标记应用使用的端口 |
| CMD | 设置默认启动命令 | 运行阶段 | 否 | 指定容器启动时执行程序 |
| ENTRYPOINT | 设置固定入口程序 | 运行阶段 | 否 | 固定容器启动主程序 |
| VOLUME | 声明数据卷 | 运行阶段 | 否 | 数据持久化 |
| USER | 指定运行用户 | 运行阶段 | 否 | 非 root 用户运行程序 |
| LABEL | 添加镜像元数据 | 构建阶段 | 否 | 添加作者、版本等信息 |
| HEALTHCHECK | 健康检查 | 运行阶段 | 否 | 判断服务是否正常 |
| ONBUILD | 设置镜像被继承时执行命令 | 构建阶段 | 否 | 制作基础镜像 |