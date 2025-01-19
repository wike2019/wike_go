package ctl

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"strings"
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

		if searchTag == "true" && !value.IsZero() {
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf(" 	 `%s` %s ?", columnName, op), value.Interface())
		}
		if searchTag == "or" && !value.IsZero() {
			columnMap[fieldName] = columnName
			db = db.Or(fmt.Sprintf(" `%s` %s ?", columnName, op), value.Interface())
		}

		if searchTag == "like" && !value.IsZero() {
			columnMap[fieldName] = columnName
			db = db.Where(fmt.Sprintf(" `%s` like ?", columnName), fmt.Sprintf("%%%v%%", value.Interface()))
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

func CreateItem[T any](db *gorm.DB, obj T) error {
	return db.Create(obj).Error
}
func UpdateItem[T any](db *gorm.DB, obj T) error {
	return db.Save(obj).Error
}
func DeleteItem[T any](db *gorm.DB, obj T) error {
	return db.Delete(obj).Error
}

func ListItem[T any](db *gorm.DB, Offset, Count int, obj interface{}, call DBChange) ([]T, int64, error) {
	var list []T
	count := int64(0)
	db = db.Model(obj)
	db, err := GetGormColumnMap(obj, db)
	if err != nil {
		return list, 0, err
	}
	if call != nil {
		db = call(db)
	}
	err = db.Count(&count).Error
	if err != nil {
		return list, 0, err
	}

	err = db.Limit(Count).Offset(Offset).Debug().Find(&list).Error
	if err != nil {
		return list, 0, err
	}
	return list, count, nil
}
func GetItemByID[T any](db *gorm.DB, obj T, id int, call DBChange) (T, error) {
	var res T
	if call != nil {
		db = call(db)
	}
	err := db.Where("id=?", id).First(&res).Error
	return res, err
}

func GetItem[T any](db *gorm.DB, obj interface{}, call DBChange) (T, error) {
	var res T
	db = db.Model(obj)
	db, err := GetGormColumnMap(obj, db)
	if err != nil {
		return nil, err
	}
	if call != nil {
		db = call(db)
	}
	err = db.Debug().First(&res).Error

	if err != nil {
		return nil, err
	}
	return res, nil
}

func ListItemAll[T any](db *gorm.DB, obj interface{}, call DBChange) ([]T, error) {
	var list []T
	db = db.Model(obj)
	db, err := GetGormColumnMap(obj, db)
	if err != nil {
		return list, err
	}
	if call != nil {
		db = call(db)
	}
	err = db.Debug().Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}
