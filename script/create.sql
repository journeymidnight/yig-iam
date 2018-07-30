grant all on iam.* to iam@"127.0.0.1" identified by "123456" with grant option;

create table User(
userName varchar(50) NOT NULL primary key,
password varchar(20) NOT NULL,
userId varchar(20) NOT NULL,
type varchar(10) NOT NULL, 
email varchar(50) NOT NULL unique,
displayName varchar(50),
accountId varchar(20),
status varchar(10),
created varchar(50) NOT NULL ,
updated varchar(50) NOT NULL
)default charset=utf8;

create table Project(
projectId varchar(20) primary key,
projectName varchar(50),
projectType varchar(20),
accountId varchar(20),
owner varchar(20),
description varchar(50),
created varchar(50),
updated varchar(50)
)default charset=utf8;

create table UserProject(
accessKey varchar(20),
accessSecret varchar(40),
acl varchar(20),
projectId varchar(20),
userId varchar(20),
created varchar(50)
CONSTRAINT pk_up PRIMARY KEY (projectId,userId)
)default charset=utf8;

create table Token(
token varchar(40) primary key,
userName varchar(50),
userId varchar(20),
accountId varchar(50),
type varchar(10),
created varchar(50),
expired varchar(50)
)default charset=utf8;

INSERT INTO User VALUES ("root", "root", "ROOT", "root@root.com", "root", "u-root", "active", now(), now());
INSERT INTO User VALUES ("account", "account", "ACCOUNT", "account@account.com", "default-account", "u-account", "active", now(), now());
