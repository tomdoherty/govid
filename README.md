# govid

[![Go Report Card](https://goreportcard.com/badge/github.com/tomdoherty/govid)](https://goreportcard.com/report/github.com/tomdoherty/govid)
[![Go Actions Status](https://github.com/tomdoherty/govid/workflows/Go/badge.svg)](https://github.com/tomdoherty/govid/actions)

## dump covid stats into docker mysql database


```shell
docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -p 3306:3306 -d mysql:latest
brew install mysql-client

```

```shell
/usr/local/opt/mysql-client/bin/mysql -h127.0.0.1 -uroot -pmy-secret-pw -t -t<<EOF
create database covid;
use covid;
create table covid (DateRep date, Cases int, Deaths int, GeoID varchar(100), Population int);
EOF

```

```shell
curl -sL https://opendata.ecdc.europa.eu/covid19/casedistribution/json >covid.json

go build -o govid cmd/main.go
./govid

```
