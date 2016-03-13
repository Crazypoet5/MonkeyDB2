## MonkeyDB2 Readme

### A stable database
### 一个稳定的数据库

MonkeyDB2 use memory storage rather than hard disk or SSD, but it doesn't means data will loss easily.<br/>
MoneyDB2使用内存储存而不是硬盘或者固态硬盘，但这并不意味着数据会很轻易的丢失。<br/>
MonkeyDB2 use technologys below to ensure stable.<br/>
MonkeyDB2使用以下技术来确保数据的稳定：<br/>

* Persistent cache sensitive index
* 可持久化的缓存敏感型索引
* File image maped memory objects
* 基于磁盘文件镜像的内存对象
* Shadow pages for transaction
* 处理事务的影子页

### A high performance database
### 一个高性能的数据库

* MonkeyDB2 use memory to storage both data and index, so it will be faster when query or insert.<br/>
* MonkeyDB2使用内存来储存数据和索引，在查询和插入时会更快。
* MonkeyDB2 use cache sensitive index, when it can use index to query, there will be a high performance.<br/>
* MonkeyDB2使用缓存敏感索引，当查询用到索引时，将会有一个良好的表现。
* MonkeyDB2 use a good query optimizer which will simplify user query and get a quick execute plan.<br/>
* MonkeyDB2使用了一个不错的查询优化器，可以简化用户查询来获得一个快速的执行计划。

### A multifunctional database
### 一个多功能的数据库

* MonkeyDB2 support both standard SQL and NoSQL to process data.
* MonkeyDB2同时支持标准SQL和NoSQL来处理数据
* MonkeyDB2 provide many utils to monitor and improve running state.
* MonkeyDB2提供了大量工具来监视和改善运行状态
* MonkeyDB2 support distributed arrangement, parted table and parted database.
* MonkeyDB2支持分布式部署，分库分表
* MonkeyDB2 support both unix-like and windows
* MonkeyDB2同时支持类Unix系统(linux, Mac, ...)和Windows系统

### A security database
### 一个安全的数据库

* MonkeyDB2 will not accept hard insertion without prepare or bind.
* MonkeyDB2不接受未经过参数绑定或准备的硬插入
* MonkeyDB2 can listen on unix socket or tcp, and with a white/black list.
* MonkeyDB2可以侦听unix套接字或者TCP，设置黑白名单
* MonkeyDB2 can transport data with encoding, even separate for data and command.
* MonkeyDB2可以编码传输数据，甚至为数据和命令使用不同的编码