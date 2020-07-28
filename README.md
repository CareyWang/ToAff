# ToAff

## Usage
```shell script
docker run -d --name toaff --restart always careywong/toaff:latest -u https://remote.config.url

OR

./toaff -u https://remote.config.url
```

## Config
支持本地当前路径下配置 config.json 或使用 -u 参数引入远程 config，配置需符合 json 语法要求。新增配置无需重启服务，支持配置自动更新。

访问 gayhub.example.com 即对应访问 https://github.com/ (大小写敏感)。

```
{
    "gayhub": "https://github.com/",
    "my": "https://github.com/CareyWang/ToAff"
}
```

## Nginx配置参考
Nginx 配置域名泛解析，使用其他程序反代请自行添加 X-Forwarded-Host 请求头。

```
server {
    server_name *.example.com;
    location / {
        proxy_http_version                 1.1;
        proxy_cache_bypass                 $http_upgrade;

        # Proxy headers
        proxy_set_header Upgrade           $http_upgrade;
        proxy_set_header Connection        "upgrade";
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host  $host;
        proxy_set_header X-Forwarded-Port  $server_port;

        # Proxy timeouts
        proxy_connect_timeout              60s;
        proxy_send_timeout                 60s;
        proxy_read_timeout                 60s;

        proxy_pass http://127.0.0.1:8006/;
    }
}
```
