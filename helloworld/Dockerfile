#基础镜像
FROM alpine:latest 
#作者信息
MAINTAINER lituquan 1242441055@qq.com 
#暴露端口
EXPOSE 8080 
#拷贝编译后的工程
ADD hello /var/paas/project-hello/
#授权可执行 
RUN chmod -R 750 /var/paas/ 
#启动工程
CMD ["/var/paas/project-hello/hello"] 
