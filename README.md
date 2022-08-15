# 项目说明
通过 redis 简单记录排行榜，
# 运行方式
本地依赖：6379 端口运行 redis server
# 接口 curl
- 更新得分，可替换 userid 和 score
```shell
curl --location --request GET '0.0.0.0:8080/ranking/updatemyscore?userid=abcddfless&score=85'
```
- 查询得分，可替换 userid
```shell
curl --location --request GET '0.0.0.0:8080/ranking/getmyranking?userid=abccc'
```
# 可能的瓶颈
- 数据没落盘，可通过添加 MySQL 数据库解决
- 数据量承载，10万日活理论上不存在瓶颈，月活超过百万级之后可能有查询瓶颈