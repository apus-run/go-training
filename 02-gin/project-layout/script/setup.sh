# 代码提交检查初始化

SOURCE_COMMIT=.github/pre-commit
TARGET_COMMIT=.git/hooks/pre-commit
SOURCE_PUSH=.github/pre-push
TARGET_PUSH=.git/hooks/pre-push

# copy pre-commit file if not exist.
if [ ! -f $TARGET_COMMIT ]; then
    echo "设置 git pre-commit hooks..."
    cp $SOURCE_COMMIT $TARGET_COMMIT
fi

# copy pre-push file if not exist.
if [ ! -f $TARGET_PUSH ]; then
    echo "设置 git pre-push hooks..."
    cp $SOURCE_PUSH $TARGET_PUSH
fi

# add permission to TARGET_PUSH and TARGET_COMMIT file.
test -x $TARGET_PUSH || chmod +x $TARGET_PUSH
test -x $TARGET_COMMIT || chmod +x $TARGET_COMMIT

echo "安装 golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

echo "安装 goimports..."
go install golang.org/x/tools/cmd/goimports@latest