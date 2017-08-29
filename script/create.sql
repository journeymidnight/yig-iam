grant all on iam.* to iam@"127.0.0.1" identified by "123456" with grant option;

create table User(
userName varchar(50) NOT NULL primary key,
password varchar(20) NOT NULL, 
type varchar(20) NOT NULL, 
email varchar(50) NOT NULL unique,
displayName varchar(50) ,
accountId varchar(50),
status varchar(10),
created varchar(50) NOT NULL ,
updated varchar(50) NOT NULL
)default charset=utf8;

create table Project(
projectId varchar(20) primary key,
projectName varchar(50),
accountId varchar(50),
description varchar(50),
status varchar(10),
created varchar(50),
updated varchar(50)
)default charset=utf8;

create table UserProject(
id int auto_increment primary key,
projectId varchar(20),
userName varchar(50),
created varchar(50)
)default charset=utf8;

create table ProjectService(
id int auto_increment primary key,
projectId varchar(20),
service varchar(20),
accountId varchar(50),
created varchar(50)
)default charset=utf8;

create table AkSk(
accessKey varchar(20) unique,
accessSecret varchar(40),
projectId varchar(20),
accountId varchar(50),
name varchar(20),
created varchar(50),
description varchar(50), 
primary key (accountId, name)
)default charset=utf8;

create table Token(
token varchar(40) primary key,
userName varchar(50),
accountId varchar(50),
type varchar(20),
created varchar(50),
expired varchar(50)
)default charset=utf8;


create table Region(
regionId varchar(20) primary key,
regionName varchar(50),
status varchar(10),
created varchar(50),
updated varchar(50)
)default charset=utf8;

create table Service(
serviceId varchar(20) primary key,
created varchar(50),
updated varchar(50),
regionId varchar(20),
endpoint varchar(254)
)default charset=utf8;


create table ProjectUser(
id int(11) primary key not null AUTO_INCREMENT,
user_id varchar(50),
project_id varchar(50),
project_name varchar(50),
role int(11),
created varchar(254)
)default charset=utf8;

INSERT INTO User VALUES ("root", "admin", "ROOT", "root@root.com", "root", "u-root", "active", now(), now());
