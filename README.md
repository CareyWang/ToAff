# ToAff

## Usage
```shell script
docker run -d --name toaff --restart always careywong/toaff:latest -u https://remote.config.url

OR

./toaff -u https://remote.config.url
```

## Nginx配置参考
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
