- [redis五种数据结构](#redis五种数据结构)
- [Building](#building)
- [Running](#running)

### redis五种数据结构
- 字符串 (String):
  使用 SET key value 命令设置字符串类型的键值对。

  例如：SET mykey "Hello"
- 列表 (List):
使用 LPUSH key value1 value2 或 RPUSH key value1 value2 

命令在列表左侧或右侧添加元素。

例如：LPUSH mylist "World"
- 集合 (Set):
使用 SADD key member1 member2 命令向集合添加成员。

例如：SADD myset "one"
- 哈希 (Hash):
使用 HSET key field1 value1 field2 value2 命令为哈希表设置字段和值。

例如：HSET myhash field1 "value1"
- 有序集合 (Sorted Set):
使用 ZADD key score1 member1 score2 member2 命令向有序集合添加成员，每个成员都有一个分数。

例如：ZADD myzset 1 "member1"

test data 
```
localhost:examples zhangzheming$ redis-cli 
127.0.0.1:6379> set apple apple
OK
127.0.0.1:6379> set boy boy
OK
127.0.0.1:6379> set test testtesttest
OK
127.0.0.1:6379> lpush book helloworld
(integer) 1
127.0.0.1:6379> lpush book hedbook
(integer) 2
127.0.0.1:6379> lpush humman man
(integer) 1
127.0.0.1:6379> lpush humman woman
(integer) 2
127.0.0.1:6379> lpush humman unknown
(integer) 3
127.0.0.1:6379> sadd myset one
(integer) 1
127.0.0.1:6379> sadd myset two
(integer) 1
127.0.0.1:6379> sadd yourset one
(integer) 1
127.0.0.1:6379> sadd yourset two
(integer) 1
127.0.0.1:6379> sadd yourset three
(integer) 1
127.0.0.1:6379> hset myhash field1 value
(integer) 1
127.0.0.1:6379> zadd myzset 1 "member1"
(integer) 1
127.0.0.1:6379> zadd myzset 2 "member2"
(integer) 1
```

### Building
``````
cd examples
go build -o analysis 
``````

### Running 
(redis 127.0.0.1:6379)
./analysis 

输出redis五种数据类型的json
```
{"hashType":[{"Key":"myhash","Type":"hash","Size":76}],"listType":[{"Key":"humman","Type":"list","Size":154},{"Key":"book","Type":"list","Size":152}],"setType":[{"Key":"yourset","Type":"set","Size":268},{"Key":"myset","Type":"set","Size":235}],"strType":[{"Key":"test","Type":"string","Size":62},{"Key":"apple","Type":"string","Size":56},{"Key":"boy","Type":"string","Size":52}],"totalType":[{"Key":"yourset","Type":"set","Size":268},{"Key":"myset","Type":"set","Size":235},{"Key":"humman","Type":"list","Size":154},{"Key":"book","Type":"list","Size":152},{"Key":"myzset","Type":"zset","Size":83},{"Key":"myhash","Type":"hash","Size":76},{"Key":"test","Type":"string","Size":62},{"Key":"apple","Type":"string","Size":56},{"Key":"boy","Type":"string","Size":52}],"zsetType":[{"Key":"myzset","Type":"zset","Size":83}]}
```

