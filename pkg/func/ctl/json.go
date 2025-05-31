package ctl

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Arr []interface{}

// MarshalJSON 自定义 JSON 序列化
func (t Arr) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}(t))
}

// UnmarshalJSON 自定义 JSON 反序列化
func (t *Arr) UnmarshalJSON(data []byte) error {
	var m []interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*t = m
	return nil
}

// Value 将 Json 转换为数据库存储的字符串
func (t Arr) Value() (driver.Value, error) {
	if t == nil {
		return "[]", nil
	}
	return json.Marshal(t)
}

// Scan 从数据库中加载 JSON 数据
func (t *Arr) Scan(value interface{}) error {
	switch v := value.(type) {
	case string: // 如果是字符串
		return t.parseJsonString(v)
	case []byte: // 如果是字节数组
		return t.parseJsonString(string(v))
	default: // 其他类型无法处理
		return errors.New("failed to scan Json: value is not a valid JSON type")
	}
}

// 解析 JSON 字符串为 []interface{}
func (t *Arr) parseJsonString(str string) error {
	var m []interface{}
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		return err
	}
	*t = m
	return nil
}

// Json 定义为 map[string]interface{} 的别名
type Json map[string]interface{}

// MarshalJSON 自定义 JSON 序列化
func (t Json) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(t))
}

// UnmarshalJSON 自定义 JSON 反序列化
func (t *Json) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*t = m
	return nil
}

// Value 将 Json 转换为数据库存储的字符串
func (t Json) Value() (driver.Value, error) {
	if t == nil {
		return "{}", nil
	}
	return json.Marshal(t)
}

// Scan 从数据库中加载 JSON 数据
func (t *Json) Scan(value interface{}) error {
	switch v := value.(type) {
	case string: // 如果是字符串
		return t.parseJsonString(v)
	case []byte: // 如果是字节数组
		return t.parseJsonString(string(v))
	default: // 其他类型无法处理
		return errors.New("failed to scan Json: value is not a valid JSON type")
	}
}

// 解析 JSON 字符串为 map[string]interface{}
func (t *Json) parseJsonString(str string) error {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		return err
	}
	*t = m
	return nil
}
