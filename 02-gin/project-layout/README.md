# 运行服务
    
```bash 
  $ go run main.go -conf=../../config 
  # 或者 
  $ CONFIG_PATH=../../config go run main.go
```

## Docker Compose 基本命令

```bash
docker compose up -d # 启动服务, -d 进程后台挂起
docker compose up # 启动服务
docker compose ps # 查看服务状态
docker compose logs -f # 查看服务日志
docker compose down # 停止服务
```

## 规范
系统自身视角的三种方式：
1、普通风格：增-Add、删-Delete、改-Update、查-Get
2、Restful 风格：增-Post、删-Delete、改-Put、查-Get
3、DB 风格：增-Insert、删-Delete、改-Update、查-Find
