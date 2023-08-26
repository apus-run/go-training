## 性能测试

可以用 apt install wrk 或者 Mac 上 brew install wrk。
源码安装：直接源码下载 git clone https://github.com/wg/wrk.git
而后进去这个 wrk 目录下，执行 make 命令编译。
编译之后你会得到一个 wrk 可执行文件，将它加入你的环境变量。

### 使用 wrk 进行压测
安装 wrk
```bash
sudo apt-get install wrk
# or
brew install wrk
```

```bash
wrk -t1 -c200 -d1s -s ./scripts/wrk/login.lua http://localhost:8080/login

wrk -t1 -c200 -d1s -s ./scripts/wrk/signup.lua http://localhost:8080/signup

wrk -t1 -c2 -d1s -s ./scripts/wrk/profile.lua http://localhost:8080/profile
```

参数说明：
```text
-t：后面跟着的是线程数量。
-d：后面跟着的是持续时间，比如说 1s 是一秒，也可以是 1m，是一分钟。
-c：后面跟着的是并发数。
-s：后面跟着的是测试的脚本。
```
