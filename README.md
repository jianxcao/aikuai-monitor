# 爱快监控

同时可以监控多台爱快的流量






```
version: '3.3'
services:
  aikuai-monitor:
    container_name: aikuai-monitor
    image: jianxcao/aikuai-monitor:v1.0.0-0-ga9b1a5-dirty
    volumes:
      - /share/docker/aikuai-monitor:/app/config
    ports:
      - 7575:7575
    environment:
      - PUID=1000
      - PGID=100
      - UMASK=0000
      - TZ=Asia/Shanghai
    restart: unless-stopped

```


/app/config为配置文件

配置文件的字段

```
[爱快路由器名称]
user=爱快登录用户名
password=爱快登录密码
url=访问爱快的地址
```