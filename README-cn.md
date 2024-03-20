[English](README.md) | 简体中文

# Hipush

Hipush是一个推送服务器，集成多个移动平台的推送通知，它支持HTTP和GRPC接口。

## Configuration

See the default [YAML config example](example.yaml):
```yaml
http:
  enabled: true
  address: "0.0.0.0"
  port: 7070

grpc:
  enabled: true
  address: "0.0.0.0"
  port: 7071

# 数据持久化配置
storage:
  enabled: true
  # 存储类型 memory redis
  type: "memory"
  # 本地持久化路径 默认路径 /etc/hipush/data.json
  path: ""

# Apns官方文档，以获取APNs集成所需的配置参数或者其他说明。
# https://developer.apple.com/documentation/usernotifications/setting-up-a-remote-notification-server
ios:
  - enabled: true               # 是否启用com.hitosea.test1应用的推送
    # 应用程序的 Bundle ID
    # ios capacitor.config文件中的appId 例如com.hitosea.test1
    app_id: "com.hitosea.test1"
    # 应用名称 方便推送的时候指定应用
    # 推送请求指定应用时可以使用app_id或者app_name指定
    app_name: "cossim"          
    key_path: ""                # APNs 密钥文件路径
    key_type: pem               # 密钥类型（例如：pem）
    password: ""                # 密钥文件的密码（如果有）
    max_concurrent_pushes: 100  # 最大并发推送数
    max_retry: 5                # 默认最大重试次数
    key_id: ""                  # 密钥 ID
    team_id: ""                 # 开发团队 ID

  - enabled: true               # 是否启用com.hitosea.test2应用的推送
    app_id: "com.hitosea.test2"
    key_path: ""
    key_type: pem
    password: ""
    production: false
    max_concurrent_pushes: 100
    max_retry: 0
    key_id: ""
    team_id: ""

# FCM官方文档，以获取FCM集成所需的配置参数或者其他说明。
# https://firebase.google.com/docs/admin/setup?hl=zh&authuser=0&_gl=1*71lbme*_ga*MjA2NzIzODYzMy4xNzA5NzE5OTQy*_ga_CW55HF8NVT*MTcxMDg0MDM4OC4yLjEuMTcxMDg0MDYwMC40Ny4wLjA.
android:
  - enabled: true
    app_id: ""
    app_name: ""
    key_path: ""        # Firebase Admin SDK AccountKey.json
    max_retry: 0        # 默认最大重试次数

# https://developer.huawei.com/consumer/cn/doc/HMSCore-References/overal-description-support-0000001064490273
huawei:
  - enabled: true
    app_id: "huawei-appid-1"
    app_secret: "huawei-app-secret-1"
    max_retry: 5

# https://dev.vivo.com.cn/documentCenter/doc/541
vivo:
  - enabled: true
    app_id: ""
    app_key: ""
    app_secret: ""
    max_retry: 5

# https://open.oppomobile.com/new/developmentDoc/info?id=10195
oppo:
  - enabled: true
    app_id: ""
    app_key: ""
    app_secret: ""     # 这里的secret是oppo的masterSecret
    max_retry: 5

# https://dev.mi.com/distribute/doc/details?pId=1529
xiaomi:
  - enabled: true
    app_id: ""
    # 小米推送需要提供包名，支持多包名
    package:
      - xxx.xx.xx
    app_secret: ""
    max_retry: 5

# https://github.com/MEIZUPUSH/PushAPI/blob/master/README.md
meizu:
  - enabled: true
    app_id: ""
    package: ""
    app_key: ""
    max_retry: 5

# https://developer.hihonor.com/cn/kitdoc?category=%E5%9F%BA%E7%A1%80%E6%9C%8D%E5%8A%A1&kitId=11002&navigation=guides&docId=kit-history.md&token=
honor:
  - enabled: true
    app_id: ""
    client_id: ""
    client_secret: ""
    max_retry: 5
```

## Deploy

直接运行项目
```bash
go run cmd/main.go -config xxx.yaml
```

使用Docker运行项目
```bash
docker run -d --name hipush \
  -v "$(pwd)/config.yaml:/config/config.yaml" \
  -p 7070:7070 \
  -p 7071:7071 \
  hub.hitosea.com/cossim/hipush \
  -config /config/config.yaml
```

使用Docker Compose运行项目 [docker-compose.yaml](docker-compose.yaml)
```bash
docker-compose up -d 
```

## Usage

更多示例在 [example](example)

### HTTP
Pushing to iOS
```markdown
curl --location --request POST 'http://<hipush-server>:7070/api/v1/push' \
--header 'Content-Type: application/json' \
--data-raw '{
    "platform": "ios",
    "token": [
        "xxxxx",
        "xxxxx"
    ],
    "app_id": "com.hitosea.test1",
    "app_name": "cossim",
    "data": {
        "title": "cossim",
        "content": "hello",
        "badge": 1,
        "sound": {
            "critical": 1,
            "volume": 4.5,
            "name": ""
        }
    },
    "option": {
        "dry_run": false,
        "retry": 2,
        "retry_interval": 1
    }
}'
