package app

import (
	"fmt"
	"strconv"

	"github.com/MeetFootball/Scaffold/crontab"
	"github.com/MeetFootball/Scaffold/model"
)

// Writer 数据库日志输出
type Writer struct{}

func (w Writer) Printf(f string, args ...interface{}) {
	const (
		LocationIndex = iota
		CostIndex
		RowsIndex
		SQLIndex
	)
	var (
		ok, slow         bool
		cost             float64
		rows             uint64
		location, t, SQL string
	)
	Log := &model.SQLLog{}
	t = "logs"
	ErrorIndex := 0
	if len(args) == 5 {
		ErrorIndex = 1
		Log.Error = fmt.Sprintf("%s", args[ErrorIndex])
		t = "errors"
		if GetEnv() == "Local" || GetEnv() == "Dev" {
			fmt.Println(args...)
		}
	}
	if location, ok = args[LocationIndex].(string); ok {
		Log.Location = location
	}
	if cost, ok = args[CostIndex+ErrorIndex].(float64); ok {
		Log.Cost = cost
		if cost > 200 {
			slow = true
			t = "slow"
		}
		Log.Slow = slow
	}
	if rows, ok = args[RowsIndex+ErrorIndex].(uint64); ok {
		Log.Rows = rows
	}
	if SQL, ok = args[SQLIndex+ErrorIndex].(string); ok {
		Log.SQL = SQL
	}
	Log.Type = t
	Log.Tag = map[string]string{"slow": strconv.FormatBool(slow), "type": t}
	crontab.SQLLogChan <- Log
}
