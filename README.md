English | [简体中文](README-cn.md)


# Hipush

Hipush is a push server that integrates push notifications for multiple mobile platforms and supports HTTP and GRPC interfaces.
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

# Data Persistence Configuration
storage:
  enabled: true
  # Storage type: memory or redis
  type: "memory"
  # Local persistence path, default path: /etc/hipush/data.json
  path: ""

# The link directs users to Apns official documentation for obtaining the required configuration parameters for APNs integration.
# https://developer.apple.com/documentation/usernotifications/setting-up-a-remote-notification-server
ios:
  - enabled: true               # Whether to enable push for com.hitosea.test1 application
    # Bundle ID of the application
    # ios capacitor.config file's appId, for example: com.hitosea.test1
    app_id: "com.hitosea.test1"
    # Application name for easy identification during push
    # Either app_id or app_name can be used to specify the application during push requests
    app_name: "cossim"
    key_path: ""                # APNs certificate file path
    key_type: pem               # Certificate type (e.g., pem)
    password: ""                # Certificate password (if any)
    max_concurrent_pushes: 100  # Maximum concurrent pushes
    max_retry: 5                # Default maximum retry attempts
    key_id: ""                  # Key ID
    team_id: ""                 # Team ID

  - enabled: true               # Whether to enable push for com.hitosea.test2 application
    app_id: "com.hitosea.test2"
    key_path: ""
    key_type: pem
    password: ""
    production: false
    max_concurrent_pushes: 100
    max_retry: 0
    key_id: ""
    team_id: ""

# The link directs users to FCM official documentation for obtaining the required configuration parameters for FCM integration.
# https://firebase.google.com/docs/admin/setup?hl=zh&authuser=0&_gl=1*71lbme*_ga*MjA2NzIzODYzMy4xNzA5NzE5OTQy*_ga_CW55HF8NVT*MTcxMDg0MDM4OC4yLjEuMTcxMDg0MDYwMC40Ny4wLjA.
android:
  - enabled: false
    app_id: ""
    app_name: ""
    key_path: ""        # Firebase Admin SDK AccountKey.json
    max_retry: 0        # Default maximum retry attempts

# https://developer.huawei.com/consumer/cn/doc/HMSCore-References/overal-description-support-0000001064490273
huawei:
  - enabled: false
    app_id: "huawei-appid-1"
    app_secret: "huawei-app-secret-1"
    max_retry: 5

# https://dev.vivo.com.cn/documentCenter/doc/541
vivo:
  - enabled: false
    app_id: ""
    app_key: ""
    app_secret: ""
    max_retry: 5

# https://open.oppomobile.com/new/developmentDoc/info?id=10195
oppo:
  - enabled: false
    app_id: ""
    app_key: ""
    app_secret: ""     # Oppo's master secret
    max_retry: 5

# https://dev.mi.com/distribute/doc/details?pId=1529
xiaomi:
  - enabled: false
    app_id: ""
    # Xiaomi push requires providing package name(s), supports multiple package names
    package:
      - xxx.xx.xx
    app_secret: ""
    max_retry: 5

# https://github.com/MEIZUPUSH/PushAPI/blob/master/README.md
meizu:
  - enabled: false
    app_id: ""
    package: ""
    app_key: ""
    max_retry: 5

# https://developer.hihonor.com/cn/kitdoc?category=%E5%9F%BA%E7%A1%80%E6%9C%8D%E5%8A%A1&kitId=11002&navigation=guides&docId=kit-history.md&token=
honor:
  - enabled: false
    app_id: ""
    client_id: ""
    client_secret: ""
    max_retry: 5
```

## Deploy

Running the Project Directly
```bash
go run cmd/main.go -config xxx.yaml
```

Running the Project Using Docker
```bash
docker run -d --name hipush \
  -v "$(pwd)/config.yaml:/config/config.yaml" \
  -p 7070:7070 \
  -p 7071:7071 \
  hub.hitosea.com/cossim/hipush \
  -config /config/config.yaml
```

Running the Project Using Docker Compose [docker-compose.yaml](docker-compose.yaml)
```bash
docker-compose up -d 
```

## Usage

More examples in [example](example)

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
