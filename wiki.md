wiki
========
>[Mysql数据库相关](https://dev.mysql.com/doc/refman/5.7/en/charset-applications.html)

1. 创建数据库 english
`CREATE DATABASE english
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;`
2. 创建单词表 word
id   name phonogram audio explanation example
`CREATE TABLE word
(
id integer primary key not null auto_increment,
name varchar(255),
phonogram varchar(255),
audio     varchar(255),
explanation varchar(512),
example    text,
createdAt  timestamp default current_timestamp,
updatedAt  timestamp default current_timestamp
);
`
3. 插入测试数据
`INSERT INTO word (name,phonogram,audio,explanation,example)
VALUES ("slogan","[ˈsloʊgən]","","n.标语，口号； 呐喊声； （商业广告上用的）短语；","Your name, logo, slogan , even the location you choose and your pricing structure depend on the brand you are trying to create.");
`

4. 安装和使用[Go Mysql Driver](http://go-database-sql.org/index.html)
`go get github.com/go-sql-driver/mysql`

[sql语句文档](https://www.w3schools.com/sql/default.asp)