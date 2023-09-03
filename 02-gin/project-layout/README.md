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

## 思考
对于 entity 是否有验证方法, 虽然推荐充血模型, 充血模型主要还是给业务层提供便利和收敛,  就不推荐加验证方法了; 
内紧外松原则, 所有的验证前置放到 handler 去验证,  保持下层的干净, 也便于测试和维护;

## 测试模板
1. 名字:简明扼要说清楚你测试的场景，建议用中文。
2. 预期输入:也就是作为你方法的输入。如果测试的是定义在类型上的方法，那么也可以包含类型实例。
3. 预期输出:你的方法执行完毕之后，预期返回的数据。 如果方法是定义在类型上的方法，那么也可以包含执行之后的实例的状态。
4. mock:每一个测试需要使用到的mock状态。单元测试 里面常见，集成测试一般没有。
5. 数据准备:每一个测试用例需要的数据。集成测试里常 见。
6. 数据清理:每一个测试用例在执行完毕之后，需要执行 一些数据清理动作。集成测试里常见。
```go
func TestXXXX(t *testing.T) {
    testCases := []struct {
        // 用例名称及描述
        name string
        
        // 预期输入, 根据你的方法参数、接收器来设计
        
        // 预期输出, 根据你的方法返回值、接收器来设计
        
        // mock 数据, 在单元测试里很常见
        mock func(ctrl *gomock.Controller)
        
        // 测试用例准备环境、数据等
        before func(t *testing.T)
        
        // 数据清理等
        after func(t *testing.T)
    }{
        {},
        {},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
        
        })
    }
}
```

## 生成 mock 文件
已经封装到 Makefile,
执行
```bash
make mock
```