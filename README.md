# Database to REST API bridge

This project aims to provide generic REST API to perform simple CRUD operations against database backend.
Currently only MySQL/MariaDB backend is supported.

### Example config

```yaml
---
server:
  http_listen_address: 0.0.0.0:22001
  cors:
    allowed_origins: ["*"]
    max_age: 1200
backends:
  demo:
    create: true # default value of `create` is false
    update: true # default value of `update` is false
    delete: true # default value of `delete` is false
    dsn: "demo:demo@tcp(localhost)/demo?parseTime=true"
  demo-ro:
    dsn: "demo:demo@tcp(localhost)/demo?parseTime=true"

```

## Security

Security is hard, so I won't even pretend :innocent:.
To secure this thing, use your favorite reverse proxy, such as nginx.

You can limit what CRUD methods are allowed.
By default, anything other than read is **NOT** allowed.

