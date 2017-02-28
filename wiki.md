wiki
========
>Mysql数据库相关
https://dev.mysql.com/doc/refman/5.7/en/charset-applications.html
1. 创建数据库 english
`CREATE DATABASE english
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;`
2. 创建单词表 word
id   name phonogram url meaning content
int  varchar()