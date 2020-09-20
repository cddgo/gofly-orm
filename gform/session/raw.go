package session

import (
	"database/sql"
	"github.com/QXQZX/gofly-orm/gform/log"
	"strings"
)

//Session 负责与数据库的交互，那交互前的准备工作（比如连接/测试数据库）
//封装有两个目的，一是统一打印日志（包括 执行的SQL 语句和错误日志）。
//二是执行完成后，清空操作。这样 Session 可以复用，开启一次会话，
//可以执行多次 SQL。
type Session struct {
	db  *sql.DB
	sql strings.Builder
	// sql 中占位符对应的值
	sqlVars []interface{}
}

func New(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, vars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, vars...)
	return s
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// package the QueryRow() method
func (s *Session) Query() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// package the Query() method
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}