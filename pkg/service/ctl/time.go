package ctl

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LocalTimeInt int64

func (t LocalTimeInt) MarshalJSON() ([]byte, error) {
	tTime := time.Unix(int64(t), 0)
	// 如果时间值是空或者0值 返回为null 如果写空字符串会报错
	if &t == nil || tTime.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

// UnmarshalJSON 为 LocalTimeInt 类型自定义 JSON 反序列化方法
func (lti *LocalTimeInt) UnmarshalJSON(data []byte) error {
	// 将 JSON 中的数字解析为 int64
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}

	// 将 int64 的时间戳转换为 time.Time，然后更新 LocalTimeInt 的值
	*lti = LocalTimeInt(timestamp)
	return nil
}

// ToTime 方法将 LocalTimeInt 类型转换为 time.Time
func (lti LocalTimeInt) ToTime() time.Time {
	return time.Unix(int64(lti), 0)
}

type LocalTime time.Time

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *LocalTime) String() string {
	// 如果时间 null 那么我们需要把返回的值进行修改
	if t == nil || t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s", time.Time(*t).Format("2006-01-02 15:04:05"))
}

func (t *LocalTime) IsZero() bool {
	return time.Time(*t).IsZero()
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")

	local, _ := time.LoadLocation("Asia/Shanghai")

	t1, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, local)
	*t = LocalTime(t1)
	return err
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	// 如果时间值是空或者0值 返回为null 如果写空字符串会报错
	if &t == nil || t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", tTime.Format("2006-01-02 15:04:05"))), nil
}
