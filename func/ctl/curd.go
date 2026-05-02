package ctl

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type DBChange func(db *gorm.DB) *gorm.DB

// GetGormColumnMap 解析结构体并返回字段名到数据库列名的映射
func GetGormColumnMap(obj interface{}, db *gorm.DB) (*gorm.DB, error) {
	columnMap := make(map[string]string)

	// 获取结构体的类型
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	// 如果是指针，获取其元素
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	// 确保传入的是结构体
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("传入的参数不是结构体")
	}

	// 遍历结构体的所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		// 跳过匿名字段
		if field.Anonymous {
			continue
		}

		// 获取字段名
		fieldName := field.Name

		// 获取 GORM 标签
		gormTag := field.Tag.Get("gorm")

		columnName := ""

		if gormTag != "" {
			// 解析标签，查找是否有 "column" 指定
			tags := strings.Split(gormTag, ";")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if strings.HasPrefix(tag, "column:") {
					columnName = strings.TrimPrefix(tag, "column:")
					break
				}
			}
		}

		if columnName == "" {
			// 如果没有指定 column 标签，使用默认的命名规则
			columnName = CamelCaseToSnakeCase(fieldName)
		}

		//columnMap[fieldName] = columnName
		orderTag := field.Tag.Get("order")
		if orderTag != "" {
			db = db.Order(fmt.Sprintf(" 	 `%s` %s ", columnName, orderTag))
		}

		searchTag := field.Tag.Get("search")

		op := field.Tag.Get("op")
		if op == "" {
			op = "="
		}

		if value.IsZero() {
			continue
		}

		switch searchTag {
		case "true":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` %s ?", columnName, op), value.Interface())
		case "or":
			columnMap[fieldName] = columnName
			db = db.Or(fmt.Sprintf("`%s` %s ?", columnName, op), value.Interface())
		case "like":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", columnName), fmt.Sprintf("%%%v%%", value.Interface()))
		case "leftlike":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", columnName), fmt.Sprintf("%%%v", value.Interface()))
		case "rightlike":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", columnName), fmt.Sprintf("%v%%", value.Interface()))
		case "eq":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` = ?", columnName), value.Interface())
		case "gt":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` > ?", columnName), value.Interface())
		case "gte":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` >= ?", columnName), value.Interface())
		case "lt":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` < ?", columnName), value.Interface())
		case "lte":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` <= ?", columnName), value.Interface())
		case "ne":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` != ?", columnName), value.Interface())
		case "in":
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf("`%s` IN ?", columnName), value.Interface())
		}
	}

	return db, nil
}

// CamelCaseToSnakeCase 将驼峰式命名转换为下划线式命名
func CamelCaseToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
