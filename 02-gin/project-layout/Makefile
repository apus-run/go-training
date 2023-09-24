# 初始化环境
.PHONY: setup
setup:
	@sh ./script/setup.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	@sh ./script/fmt.sh

.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: check
check:
	@$(MAKE) --no-print-directory fmt
	@$(MAKE) --no-print-directory tidy

# 生成 mock 文件
.PHONY: mock
mock:
	@mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
	@mockgen -source=./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go
	@mockgen -source=./internal/service/sms/async/service.go -package=svcmocks -destination=./internal/service/mocks/sms.mock.go
	@mockgen -source=./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go
	@mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
	@mockgen -source=./internal/repository/sms.go -package=repomocks -destination=./internal/repository/mocks/sms.mock.go
	@mockgen -source=./internal/repository/dao/user.go -package=daomocks -destination=./internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./internal/repository/dao/sms.go -package=daomocks -destination=./internal/repository/dao/mocks/sms.mock.go
	@mockgen -source=./internal/repository/cache/user/user.go -package=cachemocks -destination=./internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./internal/repository/cache/code/code.go -package=cachemocks -destination=./internal/repository/cache/mocks/code.mock.go
	@mockgen -package=redismocks -destination=./internal/repository/cache/redismocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable
	@mockgen -source=./internal/repository/cache/code/local.go -package=localmocks -destination=./internal/repository/cache/localmocks/local.mock.go

	@mockgen -source=./pkg/jwtx/types.go -package=jwtxmocks -destination=./pkg/jwtx/mocks/jwt.mock.go
	@mockgen -source=./internal/web/handler/jwt/types.go -package=jwtmocks -destination=./internal/web/handler/jwt/mocks/handler.mock.go

	@mockgen -source=./pkg/ratelimit/types.go -package=limitmocks -destination=./pkg/ratelimit/mocks/limit.mock.go

	@mockgen -source=./internal/service/sms/types.go -package=smsmocks -destination=./internal/service/sms/mocks/svc.mock.go



	@$(MAKE) tidy
	@echo "go mock finished"

# 单元测试
.PHONY: tests
	@go test -race ./..
	@echo "go tests finished"