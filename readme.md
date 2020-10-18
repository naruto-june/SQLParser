### 摘要
<p>
&nbsp;&nbsp;&nbsp;&nbsp;随着业务的不断扩张，企业的技术架构中会通过增加数据中间件来简化或统一对数据库的访问。数据中间件提供的一项常见的服务就是应用层通过传入json来获取数据库中的数据，这样应用层就需要了解json如何转化为sql语句的，以便确定设定传参，SQLParser就是为简化这一个过程而设计。因不同企业json与sql间转换协议并不相同，所以业务层可通过SQLParser将sql转换为通用json，然后再将通用json转换为业务所需要的json格式，毕竟一种json转另外一种json比较简单。
</p>
<p>
&nbsp;&nbsp;&nbsp;&nbsp;SQLParser是为了select语句而设计的，目前它支持字段，函数字段，子查询，orderby，groupby，limit等，但不支持having，union，join等关键字的sql语句，可满足业务的绝大部分需求。
</p>

### 使用
##### 编译 parser
```
    cd $GOPATH/src/SQLParser/cmd
    go build -o parser parser.go
```
##### 使用 parser
```
    cd $GOPATH/src/SQLParser/cmd
    ./parser -s "select id from table_name where left(id,2)=\"te\""
```

##### 编译 parserSrv
```
    cd $GOPATH/src/SQLParser/cmd
    go build -o parserSrv parserSrv.go
```
##### 使用 parserSrv
```
    cd $GOPATH/src/SQLParser/cmd
    ./parserSrv >>nohup.out 2>&1 &
    curl -i "localhost:8000/sql/util/parser?" -d'select id from table_name where left(id,2)="te"'
```

##### 编译并安装 selecter库
```
    cd $GOPATH/src/SQLParser/selecter
    go install
```


### 联系
```
1643127918@qq.com
```