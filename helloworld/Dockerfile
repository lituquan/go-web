FROM ubuntu:latest  #基础镜像
MAINTAINER lituquan 1242441055@qq.com #作者信息
EXPOSE 8080 #暴露端口
ADD hello /var/paas/project-hello/ #拷贝编译后的工程
RUN chmod -R 750 /var/paas/ #授权可执行
CMD ["/var/paas/project-hello/hello"] #启动镜像
