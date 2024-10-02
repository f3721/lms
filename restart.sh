#!/bin/bash
echo "go build"
#go mod tidy
go build -mod=vendor -o go-admin main.go
go build -mod=vendor -o proxy cmd/proxy/proxy.go
chmod +x ./go-admin
echo "kill go-admin service"
killall go-admin # kill go-admin service
echo "kill proxy service"
killall proxy # kill proxy service
nohup ./go-admin server -c=config/settings.admin.dev.yml >> access.admin.log 2>&1 &   #后台启动服务将日志写入access.admin.log文件
nohup ./go-admin server -c=config/settings.uc.dev.yml >> access.uc.log 2>&1 &         #后台启动服务将日志写入access.uc.log文件
nohup ./go-admin server -c=config/settings.oc.dev.yml >> access.oc.log 2>&1 &         #后台启动服务将日志写入access.oc.log文件
nohup ./go-admin server -c=config/settings.pc.dev.yml >> access.pc.log 2>&1 &         #后台启动服务将日志写入access.pc.log文件
nohup ./go-admin server -c=config/settings.wc.dev.yml >> access.wc.log 2>&1 &         #后台启动服务将日志写入access.wc.log文件
nohup ./proxy >> proxy.log 2>&1 &                                                     #后台启动服务将日志写入proxy.log文件
echo "run go-admin success"
echo "services list:"
echo "    admin port:8000"
echo "    uc port:8001"
echo "    pc port:8002"
echo "    wc port:8003"
echo "    oc port:8004"
echo "    proxy port:8888"
ps -aux | grep go-admin
ps -aux | grep proxy