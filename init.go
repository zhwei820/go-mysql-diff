package main

import (
	"context"
	"fmt"

	"github.com/zhwei820/log"
)

// TriggerDiff 对比触发器的不同
// func TriggerDiff(db1, db2 *sql.DB, schema1, schema2 string) bool {
// 	triggerName1, err := getTriggerName(db1, schema1)
// 	if err != nil {
// 		dLog.Fatalln(err.Error())
// 	}
// 	triggerName2, err := getTriggerName(db2, schema2)
// 	if err != nil {
// 		dLog.Fatalln(err.Error())
// 	}
// 	if !isEqual(triggerName1, triggerName2) {
// 		dt := diffName(triggerName1, triggerName2)
// 		dLog.Printf("两个数据库不同的触发器,共有%d个，分别是：%s", len(dt), dt)
// 		return false
// 	}
// 	// dLog.Printf("两个数据库触发器相同")
// 	return true
// }

// func getTriggerName(s *sql.DB, schema string) (ts []string, err error) {
// 	stm, perr := s.Prepare("select TRIGGER_NAME from information_schema.triggers where TRIGGER_SCHEMA=? order by TRIGGER_NAME")
// 	if perr != nil {
// 		err = perr
// 		return
// 	}
// 	defer stm.Close()
// 	q, qerr := stm.Query(schema)
// 	if qerr != nil {
// 		err = qerr
// 		return
// 	}
// 	defer q.Close()

// 	for q.Next() {
// 		var name string
// 		if err := q.Scan(&name); err != nil {
// 			log.Fatal(err)
// 		}
// 		ts = append(ts, name)
// 	}
// 	return
// }

// // FunctionDiff 对比函数的不同
// func FunctionDiff(db1, db2 *sql.DB, schema1, schema2 string) bool {
// 	functionName1, err := getFunctionName(db1, schema1)
// 	if err != nil {
// 		dLog.Fatalln(err.Error())
// 	}
// 	functionName2, err := getFunctionName(db2, schema2)
// 	if err != nil {
// 		dLog.Fatalln(err.Error())
// 	}
// 	fmt.Println(functionName1)
// 	fmt.Println(functionName2)
// 	if !isEqual(functionName1, functionName2) {
// 		dt := diffName(functionName1, functionName2)
// 		dLog.Printf("两个数据库不同的函数,共有%d个，分别是：%s", len(dt), dt)
// 		return false
// 	}
// 	// dLog.Printf("两个数据库函数相同")
// 	return true
// }

// func getFunctionName(s *sql.DB, schema string) (ts []string, err error) {
// 	stm, perr := s.Prepare("select ROUTINE_NAME from information_schema.routines where ROUTINE_SCHEMA=? and ROUTINE_TYPE='FUNCTION' order by ROUTINE_NAME")
// 	if perr != nil {
// 		err = perr
// 		return
// 	}
// 	defer stm.Close()
// 	q, qerr := stm.Query(schema)
// 	if qerr != nil {
// 		err = qerr
// 		return
// 	}
// 	defer q.Close()

// 	for q.Next() {
// 		var name string
// 		if err := q.Scan(&name); err != nil {
// 			log.Fatal(err)
// 		}
// 		ts = append(ts, name)
// 	}
// 	return
// }

func Info(msg string, arg ...interface{}) {
	log.InfoZ(context.Background(), fmt.Sprintf(msg, arg...))
}
