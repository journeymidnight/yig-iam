# 关于yig-iam
yig-iam是与yig项目配套使用的用户管理系统，负责组织用户账户和项目的对应关系，以及每个项目的ak/sk的管理

# 账户等级划分
1. root
2. account
3. user

root是根用户，用于管理account账户的增删改查
account称为账户，用于管理名下创建的user(用户)、project(项目)，并建立user-project的映射关系，每对映射关系都关联一组ak/sk
user称为普通用户，可以查看和它关联的project的相关信息(包含ak/sk)

# 与yig的逻辑对应关系

yig-iam里面的project就对应yig里面的一个用户,包含一个独立的对象存储空间

# 部署步骤

## 1.编译

```
cd ../yig-iam
make
```

## 2.打rpm包

```
cd ../yig-iam
sh package/rpmbuild.sh
```

## 3.修改yig-iam配置文件

```
mkdir /etc/yig-iam
cp .../yig-iam/config/* /etc/yig-iam/
cat /etc/yig-iam/conf.toml
```
典型配置

```
ManageKey = "key"
ManageSecret = "secret"
Logpath = "/var/log/yig-iam/iam.log"
Loglevel = "info"
Accesslog = "/var/log/yig-iam/access.log"
PidFile = "/tmp/iam.pid"
BindPort =  8888
RbacDataSource = "root:@tcp(127.0.0.1:4000)/"
UserDataSource = "root:@tcp(127.0.0.1:4000)/"
TokenExpire = 36000
```


其中：
```ManageKey```、```ManageSecret```代表yig-iam和yig进程的管理秘钥，yig侧的配置文件(```/etc/yig/yig.toml```)与yig-iam的此两项配置项必须一致;
```LogPath```、```PanicLogPath```、```PidFile```为日志和pid的路径，没必要修改；
```BindPort```为iam的绑定端口，默认为8888；
```RbacDataSource ```、```UserDataSource ```为访问数据库的路径;
```LogLeval```为日志的等级;
```TokenExpire```为login操作获取的token失效时间，单位为秒；

PS：确保数据库连接配置正确,默认连接本机mysql数据库，端口4000，用户名root，密码为空

## 4.修改yig的配置文件
```
vim /etc/yig/yig.toml
```
添加以下内容

```
[plugins.yig_iam]
path = "/etc/yig/plugins/yig_iam_plugin.so"
enable = true
[plugins.yig_iam.args]
EndPoint="http://127.0.0.1:8888/api/v1/yig/fetchsecretkey"
ManageKey="key"
ManageSecret="secret"
```

## 5.启动yig-iam

```
./yig-iam

```
初次启动后自动创建相关数据库的表，同时会生成两个默认账户，分别是root账户和account账户，默认用户名和密码分别为：
root：root
admin：admin

## 6.使用命令行工具yig-iam-tools创建用户
1.用默认admin用户登录,获取token

```
[root@node1 yig-iam]# ./yig-iam-tools login -u admin -p admin
command: login user: admin password: admin
endPoint: http://127.0.0.1:8888 token:
{
	"token": "c42cd174-f990-443e-8ef7-a048737abb9d",
	"type": "ACCOUNT",
	"userId": "u-yZ51d24Ejuzuy2fS",
	"accountId": "u-yZ51d24Ejuzuy2fS"
}
```
2.创建普通用户user1，密码user1

```
[root@node1 yig-iam]# ./yig-iam-tools createuser -u user1 -p user1 -t c42cd174-f990-443e-8ef7-a048737abb9d
command: createuser user: user1 password: user1
endPoint: http://127.0.0.1:8888 token: c42cd174-f990-443e-8ef7-a048737abb9d
```
3.使用普通用户user1，密码user1登录,  拿到token

```
[root@node1 yig-iam]# ./yig-iam-tools login -u user1 -p user1
command: login user: user1 password: user1
endPoint: http://127.0.0.1:8888 token:
{
	"token": "a6e87eb2-c0d9-4c6b-94a2-910a0ca458b0",
	"type": "USER",
	"userId": "u-MHVzXOcbx93LEN5K",
	"accountId": "u-yZ51d24Ejuzuy2fS"
}
```

4.查看新建user1用户的ak/sk

```
[root@node1 yig-iam]# ./yig-iam-tools listkeys -t a6e87eb2-c0d9-4c6b-94a2-910a0ca458b0
command: listkeys user:  password:
endPoint: http://127.0.0.1:8888 token: a6e87eb2-c0d9-4c6b-94a2-910a0ca458b0
[
	{
		"projectId": "p-KSsanigOAD0AOKoH",
		"name": "p-KSsanigOAD0AOKoH",
		"accessKey": "iuy3j5Q4f7gjlb7hQPdP",
		"accessSecret": "rCIwYiKBTAv85Tds1U9AZX7oTBhQSx3b8vMJwmYO",
		"acl": "RW",
		"status": "enable",
		"updated": "2020-04-09 14:37:46"
	}
]
```

5.如果忘记创建的用户名和密码，可用admin用户的token，查询到全部新建用户的账户和密码

```
[root@node1 yig-iam]# ./yig-iam-tools listusers -t c42cd174-f990-443e-8ef7-a048737abb9d
command: listusers user:  password:
endPoint: http://127.0.0.1:8888 token: c42cd174-f990-443e-8ef7-a048737abb9d
[
	{
		"userName": "user1",
		"password": "user1",
		"userId": "u-MHVzXOcbx93LEN5K",
		"type": "USER",
		"email": "",
		"displayName": "user1",
		"accountId": "u-yZ51d24Ejuzuy2fS",
		"status": "active",
		"created": "2020-04-09 14:37:46",
		"updated": "2020-04-09 14:37:46"
	}
]
```

## 7.帮助信息

```
[root@node1 yig-iam]# ./yig-iam-tools -h
command: -h user:  password:
endPoint: http://127.0.0.1:8888 token:
Usage: admin <commands> [options...]
Commands: login|createaccount|createuser|listaccounts|listkeys
               |deleteaccount|deleteuser|listusers|
Options:
 -e, --endpoint   Specify endpoint of yig-iam
 -u, --user      Specify user name to login
 -p, --password   Specify password to login
 -t, --token   Specify token to send request with
 -m, --email   Specify email
 -d, --displayname   Specify display name
```

## 8.Rest API
1.yig-iam内置了更丰富的api接口，请查阅相关源代码


## 9.TODO
1.使用简单的web页面来方便使用  
2.扩展命令行的使用（目前只来得及实现了部分接口）
