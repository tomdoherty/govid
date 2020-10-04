# govid


dump covid stats into mysql database

```
docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -p 3306:3306 -d mysql:latest
brew install mysql-client
/usr/local/opt/mysql-client/bin/mysql -h127.0.0.1 -uroot -pmy-secret-pw -t -t<<EOF
create database covid;
use covid;
create table covid (DateRep date, Cases int, Deaths int, GeoID varchar(100), Population int);
EOF
curl -sL https://opendata.ecdc.europa.eu/covid19/casedistribution/json >covid.json
go build
./govid
```
