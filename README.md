# Database to REST API bridge

This project aims to provide generic REST API to perform simple CRUD operations against database backend.
Currently only MySQL/MariaDB backend is supported.

### Example config

```yaml
---
server:
  listen_address: 0.0.0.0:22001
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

## Integration tests

Integration tests are written in [robotframework](https://robotframework.org/).
Have a database MySQL/MariaDB server ready. Test suite does not create database,
but it will populate (and drop afterward) the tables.

Here are environment variables that control how suite connect to database:

| Variable name    | Default values | Description                                                                                                      |
|------------------|----------------|------------------------------------------------------------------------------------------------------------------|
| `MYSQL_HOST`     | `localhost`    | Hostname where database server is running                                                                        |
| `MYSQL_PORT`     | `3306`         | TCP port on which DB server is listening on                                                                      |
| `MYSQL_DDL_USER` | `root`         | User that is used to create and drop tables                                                                      |
| `MYSQL_DDL_PASS` | `123456`       | Password for `MYSQL_DDL_USER` user                                                                               |
| `MYSQL_APP_USER` | `demo`         | User that application will use                                                                                   |
| `MYSQL_APP_PASS` | `123456`       | Password for `MYSQL_DDL_PASS` user                                                                               |
| `MYSQL_DB`       | `demo`         | Name of database that applciation will use. <br/>It's also name of _backend_ in configuration and API endpoints. |

To install all RF dependencies you can use `make prepare`, which will set up [virtualenv](https://docs.python.org/3/library/venv.html) for you.
To run actual suite, use `make it`
