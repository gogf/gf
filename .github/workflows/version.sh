#!/usr/bin/env bash

# 判断参数个数是否为2
if [ $# -ne 2 ]; then
    echo "参数异常，请以version.sh [目录] [版本号]的格式执行"
    echo "例如：version.sh ./contrib v1.0.0"
    exit 1
fi

# 判断第一个参数是否为目录并存在
if [ ! -d "$1" ]; then
    echo "错误：目录不存在"
    exit 1
fi

# 判断第二个参数是否以v开头
if [[ "$2" != v* ]]; then
    echo "错误：版本号不是以v开头"
    exit 1
fi

workdir=$1
newVersion=$2
echo "准备将${workdir}目录下的所有go.mod文件中的版本号替换为${newVersion}"

if [[ ${workdir} == ./contrib ]]; then
    echo "package gf" > version.go
    echo "" >> version.go
    echo "const (" >> version.go
    echo -e "\t// VERSION is the current GoFrame version." >> version.go
    echo -e "\tVERSION = \"${newVersion}\"" >> version.go
    echo ")" >> version.go
fi

# 判断文件是否存在
if [ -f "go.work" ]; then
    # 文件存在，重命名
    mv go.work go.work.${newVersion}
    echo "备份go.work文件，以免影响升级"
fi

for file in `find ${workdir} -name go.mod`; do
    goModPath=$(dirname $file)
    echo ""
    echo "processing dir: $goModPath"
    cd $goModPath
    go mod tidy
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf"
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf" | xargs -L1 go get -v 
    go mod tidy
    cd -
done

if [ -f "go.work.${newVersion}" ]; then
    # 文件存在，重命名
    mv go.work.${newVersion} go.work
    echo "恢复go.work文件"
fi