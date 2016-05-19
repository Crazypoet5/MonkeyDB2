# MonkeyDB2

## How to use

### Beigin

First, you should install and run the `monkeyd` server:  
``` 
./install.cmd (For windows) or ./install.sh (For linux)
./monkeyd [-port 2016]
```  
Then, you can use `monkey` shell to manage it:
``` ./monkey [-port 2016] ```
> Sorry, MonkeyDB2 only support local access up to now.  

Otherwise, you can use your favourite language, eg:  
nodejs:
```
var monkey = require("./monkey")
monkey.connect(function() {
    monkey.query("select * from t1", function(data) {
        data = JSON.parse(data)
        console.log(data)
        monkey.close()
    })
})    
```
or Python:
```
import monkey

monkey.Connect()
cmd = raw_input("Monkey>>")
while (cmd != "quit" and cmd != "quit;"):
    res = monkey.SendCmd(cmd)
    print res
    cmd = raw_input("Monkey>>")
monkey.Close()
```

And the commands supported are follow:

* show tables

```
Monkey>>show tables;
_______
|table|
|-----|
|a    |
|-----|
1  row(s) affected in  0 s.
```

* show table `table`

```
Monkey>>show table a;
__________________
|filed|type  |key|
|-----|------|---|
|id   |INT   |N  |
|-----|------|---|
|name |STRING|N  |
|-----|------|---|
1  row(s) affected in  0 s.
```

* create table `table` ( `filed` `type` `attributes`)

```
Monkey>>create table test (
      ->id int primary key,
      ->name string unique,
      ->gpa float);
0  row(s) affected in  0.002001 s.
```

* insert into `table`(`fileds`...) values (...), ...

```
Monkey>>insert into test(id, name, gpa) values (1, 'InsZVA', 2.9), (2, 'LowesYang', 5.0);
2  row(s) affected in  0 s.
```

* select `fields`... from `table` where ...

```
Monkey>>select * from test where id=2 or name='InsZVA';
__________________
|id|name     |gpa|
|--|---------|---|
|1 |InsZVA   |2.9|
|--|---------|---|
|2 |LowesYang|5  |
|--|---------|---|
2  row(s) affected in  0 s.
```

* update `table` set `field` = `value` where ...

```
Monkey>>update test set gpa=3.2 where id=1;
1  row(s) affected in  0.0009849 s.
```

* delete from `table` where ...

```
Monkey>>delete from test where id=1;
1  row(s) affected in  0 s.
```

* create index `index` on `table`(`field`)

```
Monkey>>create index in1 on test(gpa);
0  row(s) affected in  0.001001 s.
```

* drop index `index`

```
Monkey>>drop index in1;
0  row(s) affected in  0.0009831 s.
```

* drop table `table`

```
Monkey>>drop table test;
0  row(s) affected in  0.0020041 s.
```

> Below is nosql support

* createkv `kvTable` `keyType` `valueType`

```
Monkey>>createkv kv int string;
0  row(s) affected in  0.0019841 s.
```

* set `kvTable` `key` `value`

```
Monkey>>set kv 1 'abv';
1  row(s) affected in  0 s.
```

* get `kvTable` `key`

```
Monkey>>get kv 1;
_______
|value|
|-----|
|abv  |
|-----|
1  row(s) affected in  0 s.
```

* remove `kvTable` `key`

```
Monkey>>remove kv 1;
1  row(s) affected in  0 s.
```

> KVTable can also be produced by SQL-clause like below

* drop table `kvTable`

```
Monkey>>drop table kv;
0  row(s) affected in  0.0020001 s.
```

### Use it in programing

MonkeyDB2 support Node.js & Python (and more in future), you can import `monkey` package
, and send command above to access the monkey database, and the return of these functions
 is a json object with 3 elements: `relation`, `result` and `error` like:
 
 ```
 {"relation":[],"result":{"affectedRows":0,"usedTime":2017000},"error":null}
 ```
 
 You can produce this json in your own program.

## MonkeyDB2 Features

### A stable database
### 一个稳定的数据库

MonkeyDB2 use memory storage rather than hard disk or SSD, but it doesn't means data will loss easily.<br/>
MoneyDB2使用内存储存而不是硬盘或者固态硬盘，但这并不意味着数据会很轻易的丢失。<br/>
MonkeyDB2 use technologys below to ensure stable.<br/>
MonkeyDB2使用以下技术来确保数据的稳定：<br/>

* Persistent cache sensitive index  [finished]
* 可持久化的缓存敏感型索引            [finished]
* File image maped memory objects   [finished]
* 基于磁盘文件镜像的内存对象          [finished]
* Shadow pages for transaction
* 处理事务的影子页

### A high performance database
### 一个高性能的数据库

* MonkeyDB2 use memory to storage both data and index, so it will be faster when query or insert. [finished]<br/>
* MonkeyDB2使用内存来储存数据和索引，在查询和插入时会更快。 [finished]
* MonkeyDB2 use cache sensitive index, when it can use index to query, there will be a high performance.<br/>
* MonkeyDB2使用缓存敏感索引，当查询用到索引时，将会有一个良好的表现。
* MonkeyDB2 use a good query optimizer which will simplify user query and get a quick execute plan.<br/>
* MonkeyDB2使用了一个不错的查询优化器，可以简化用户查询来获得一个快速的执行计划。

### A multifunctional database
### 一个多功能的数据库

* MonkeyDB2 support both standard SQL and NoSQL to process data. [part-finished]
* MonkeyDB2同时支持标准SQL和NoSQL来处理数据 [part-finished]
* MonkeyDB2 provide many utils to monitor and improve running state.
* MonkeyDB2提供了大量工具来监视和改善运行状态
* MonkeyDB2 support distributed arrangement, parted table and parted database.
* MonkeyDB2支持分布式部署，分库分表
* MonkeyDB2 support both unix-like and windows [part-finished]
* MonkeyDB2同时支持类Unix系统(linux, Mac, ...)和Windows系统 [part-finished]

### A security database
### 一个安全的数据库

* MonkeyDB2 will not accept hard insertion without prepare or bind.
* MonkeyDB2不接受未经过参数绑定或准备的硬插入
* MonkeyDB2 can listen on unix socket or tcp, and with a white/black list.
* MonkeyDB2可以侦听unix套接字或者TCP，设置黑白名单
* MonkeyDB2 can transport data with encoding, even separate for data and command.
* MonkeyDB2可以编码传输数据，甚至为数据和命令使用不同的编码