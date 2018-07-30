package store

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/journeymidnight/nier/src/api/datatype"
	"github.com/journeymidnight/nier/src/helper"
	"time"
	"errors"
)



const TimeFormat = "2006-01-02T15:04:05Z07:00"

type DbStore struct{
	db *sql.DB
}

func NewDbStore() *DbStore {
	s := DbStore{}
	db, err := sql.Open("mysql", helper.Config.UserDataSource)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %v", err))
	}
	helper.Logger.Infoln("Connected to database")

	_,err = db.Exec("CREATE DATABASE IF NOT EXISTS iam")
	if err != nil {
		panic(err)
	}

	_,err = db.Exec("USE iam")
	if err != nil {
		panic(err)
	}

	_, err =  db.Exec("create table if not exists User( userName varchar(50) primary key, password varchar(20), type varchar(10))default charset=utf8")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists Snmp( address varchar(50) primary key)default charset=utf8") 
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists History( userName varchar(50), opType varchar(50), created varchar(50)) default charset=utf8;")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists Token( token varchar(40) primary key, userName varchar(50), type varchar(10), created varchar(50), expired varchar(50))default charset=utf8")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists alerts(eventid varchar(60) NOT NULL, eventname varchar(40), node varchar(40), sendtime datetime, status varchar(20), PRIMARY KEY (`eventid`)) default charset=utf8")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists Ntp (ntpAddress varchar(40) primary key) default charset=utf8")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists Snmp (snmpAddress varchar(40) primary key) default charset=utf8")
	if err != nil {
		panic(err)
	}

	//other nier could have inserted
	db.Exec("INSERT IGNORE INTO User VALUES ('admin', 'admin', 'ROOT')")
	db.Close()
	db, err = sql.Open("mysql", helper.Config.UserDataSource+"iam")
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %v", err))
	}
	s.db = db
	return &s
}

func (s *DbStore)InsertUserRecord(userName string, password string, accountType string) error {
	_, err := s.db.Exec("insert into User values( ?, ?, ? )", userName, password, accountType)
	if err != nil {
		helper.Logger.Errorln("Error add user", userName, password, accountType, err.Error())
		return err
	}
	return err
}

func (s *DbStore)RemoveUserRecord(userName string) error {
	_, err := s.db.Exec("delete from User where userName=(?)", userName)
	if err != nil {
		helper.Logger.Errorln("Error remove user", userName, err.Error())
	}
	return err
}

func (s *DbStore)ModifyUserRecord(userName string, password string, accountType string) error {
	_, err := s.db.Exec("update User set password=(?), type=(?) where userName=(?)", password, accountType, userName)
	if err != nil {
		helper.Logger.Errorln("Error remove user", userName, err.Error())
	}
	return err
}

func (s *DbStore)DescribeUserRecord(userName string) (UserRecord, error) {
	var record UserRecord
	err := s.db.QueryRow("select * from User where userName=(?)", userName).Scan(&record.UserName,
		&record.Password,
		&record.Type)
	return record, err
}

func (s *DbStore)ValidUserAndPassword(userName string, password string) (UserRecord, error) {
	var record UserRecord
	err := s.db.QueryRow("select * from User where userName=(?) and password=(?)", userName, password).Scan(
		&record.UserName,
		&record.Password,
		&record.Type)
	return record, err
}

func (s *DbStore)CheckUserExist(UserName string) (bool, error) {
	var count int
	err := s.db.QueryRow("select count(*) from User where userName=(?)", UserName).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *DbStore)InsertTokenRecord(Token string, UserName string, Type string) error {
	created := time.Now().Format(TimeFormat)
	expired := time.Now().Add(time.Duration(helper.Config.TokenExpire * 1000000000)).Format(TimeFormat)
	_, err := s.db.Exec("insert into Token values( ?, ?, ?, ?, ? )", Token, UserName, Type, created, expired)
	if err != nil {
		helper.Logger.Errorln("Error InsertTokenRecord", Token, UserName, Type, created, expired, err.Error())
	}
	return err
}

func (s *DbStore)RemoveTokenRecord(Token string) error {
	_, err := s.db.Exec("delete from Token where token=(?)", Token)
	if err != nil {
		helper.Logger.Errorln("Error remove token", Token, err.Error())
	}
	return err
}

func (s *DbStore)SearchExistedToken(userName string) (TokenRecord, error) {
	var record TokenRecord
	rows, err := s.db.Query("select * from Token")
	if err != nil {
		helper.Logger.Errorln("Error querying idle executors: ", err)
		return record, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&record.Token,
			&record.UserName,
			&record.Type,
			&record.Created,
			&record.Expired); err != nil {
			helper.Logger.Errorln("Row scan error: ", err)
			continue
		}
		if record.UserName == userName {
			expired := time.Now().Add(time.Duration(helper.Config.TokenExpire * 1000000000)).Format(TimeFormat)
			_, err := s.db.Exec("update Token set expired=(?) where token=(?)", expired, record.Token)
			if err != nil {
				helper.Logger.Errorln("Error update TokenRecord", record.Token, err.Error())
				return record, err
			}
			return record, nil
		}
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Errorln("Row error: ", err)
	}
	return record, errors.New("token not found")
}

func (s *DbStore)ListExpiredTokens() ([]TokenRecord, error) {
	var records []TokenRecord
	rows, err := s.db.Query("select * from Token")
	if err != nil {
		helper.Logger.Errorln("Error querying idle executors: ", err)
		return records, err
	}
	defer rows.Close()
	now := time.Now()
	for rows.Next() {
		var record TokenRecord
		if err := rows.Scan(
			&record.Token,
			&record.UserName,
			&record.Type,
			&record.Created,
			&record.Expired); err != nil {
			helper.Logger.Errorln("Row scan error: ", err)
			continue
		}
		t, err := time.Parse(TimeFormat, record.Expired)
		if err != nil {
			helper.Logger.Errorln("Error parse expired time: ", err)
			records = append(records, record)
			continue
		}
		if now.Sub(t) > 0 {
			records = append(records, record)
		}
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Errorln("Row error: ", err)
	}
	return records, err

}

func (s *DbStore)GetTokenRecord(Token string) (TokenRecord, error) {
	var record TokenRecord
	err := s.db.QueryRow("select * from Token where token=(?)", Token).Scan(&record.Token,
		&record.UserName,
		&record.Type,
		&record.Created,
		&record.Expired)
	return record, err
}

func (s *DbStore)ListUserRecords() ([]UserRecord, error) {
	var records []UserRecord
	rows, err := s.db.Query("select * from User")
	if err != nil {
		helper.Logger.Errorln("Error querying idle executors: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record UserRecord
		if err := rows.Scan(
			&record.UserName,
			&record.Password,
			&record.Type); err != nil {
			helper.Logger.Errorln("Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Errorln("Row error: ", err)
	}
	return records, err
}

func (s *DbStore)InsertSnmpRecord(snmpAddress string) error {
	_, err := s.db.Exec("delete from Snmp")
	if err != nil {
		helper.Logger.Errorln("clear snmp table failed", err.Error())
		return err
	}

	_, err = s.db.Exec("insert into Snmp values( ? )", snmpAddress)
	if err != nil {
		helper.Logger.Errorln("Error add snmp address", snmpAddress, err.Error())
		return err
	}
	return err
}

func (s *DbStore)InsertNtpRecord(ntpAddress string) error {
	_, err := s.db.Exec("delete from Ntp")
	if err != nil {
		helper.Logger.Errorln("clear ntp table failed", err.Error())
		return err
	}

	_, err = s.db.Exec("insert into Ntp values( ? )", ntpAddress)
	if err != nil {
		helper.Logger.Errorln("Error add ntp address", ntpAddress, err.Error())
		return err
	}
	return err
}

func (s *DbStore)GetSnmpRecord() (SnmpRecord, error) {
	var record SnmpRecord
	err := s.db.QueryRow("select * from Snmp").Scan(&record.SnmpAddress)
	return record, err
}

func (s *DbStore)GetNtpRecord() (NtpRecord, error) {
	var record NtpRecord
	err := s.db.QueryRow("select * from Ntp").Scan(&record.NtpAddress)
	return record, err
}

func (s *DbStore)GetHistoryRecord() ([]HistoryRecord, error) {
	var records []HistoryRecord = make([]HistoryRecord,0)
	rows, err := s.db.Query("select * from History")
	if err != nil {
		helper.Logger.Errorln("Error querying op history: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record HistoryRecord
		if err := rows.Scan(
			&record.UserName,
			&record.OpType,
			&record.Time); err != nil {
			helper.Logger.Errorln("Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Errorln("Row error: ", err)
	}
	return records, err
}

func (s *DbStore)InsertOpRecord(userName, opType string) error {
	created := time.Now().Format(TimeFormat)
	_, err := s.db.Exec("insert into History values( ?, ?, ? )", userName, opType, created)
	if err != nil {
		helper.Logger.Errorln("Error log operate record", userName, opType, err.Error())
	}
	return err
}

//golang get datetime type from mysql
//https://www.jianshu.com/p/444bf0fddcd7
//https://github.com/go-sql-driver/mysql#timetime-support
//I use string for mysql.DATETIME
func (s *DbStore)GetAlertRecord() ([]AlertRecord, error) {
	var alerts []AlertRecord = make([]AlertRecord, 0)
	rows, err := s.db.Query("select eventid, eventname, node, sendtime, status  from alerts order by sendtime DESC limit 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var alert AlertRecord
		rows.Scan(&alert.EventId, &alert.EventName, &alert.Node, &alert.SendTime, &alert.Status)
		alerts = append(alerts, alert)
	}
	return alerts, nil

}

func (s *DbStore)Close() {
	s.db.Close()
}


