#!/bin/bash
#
# 快速部署脚本
#

#go test -v -coverprofile=./test/testReport.out
#go tool cover -html=./test/testReport.out -o=./test/testReport.html

VERSION=$1
echo $VERSION
if [ "$VERSION" ];then
    echo "tag版本号为$VERSION"
else
    echo "请输入tab版本号";
    exit 1;
fi;

git pull origin master
git add .
git commit -m "fix:上传新版本"
git tag $VERSION
if [ $? -eq 0 ]; then
    echo "tag成功";
else
    echo "打tag失败。";
    exit 1
fi;
git push origin master
git push origin $VERSION