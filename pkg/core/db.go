package core

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type CoreDb struct {
	DB *gorm.DB
}

func InitDb() *CoreDb {
	dbsqlite, err := gorm.Open(sqlite.Open("core.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = dbsqlite.AutoMigrate(&API{})
	if err != nil {
		panic("failed to migrate database")
	}
	dbsqlite.Model(API{}).Where("1=1").Delete(nil)
	return &CoreDb{
		DB: dbsqlite,
	}
}

func Input(query interface{}, body interface{}, header interface{}) string {
	str := ""
	str += generateStructTable(query, 2, 1)
	str += "---\n"
	str += generateStructTable(body, 1, 1)
	str += "---\n"
	str += generateStructTable(header, 3, 1)
	str += "---\n"
	return str
}
func Output(body interface{}) string {
	return generateStructTable(body, 1, 1) + "---\n"
}
func (this *CoreDb) ApiTable(name string, group string, input string, output string, path string, method string) *CoreDb {
	res := &API{
		Name:   name,
		Group:  group,
		Input:  input,
		Output: output,
		Path:   path,
		Method: method,
	}
	this.DB.Create(res)
	return this
}

func generateStructTable(data interface{}, dataType int, level int) string {
	if data == nil {
		return ""
	}
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// 如果是指针，解引用
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	// 如果不是结构体，返回空
	if t.Kind() != reflect.Struct {
		return ""
	}

	md := strings.Builder{}

	md.WriteString(fmt.Sprintf("%s### %s\n\n", strings.Repeat("#", level), t.Name()))

	defaultData := "字段名 | 类型 | 标签 (json) | 描述 | 是否必填 |\n"
	if dataType == 2 {
		defaultData = "字段名 | 类型 | 标签 (form) | 描述 | 是否必填 |\n"
	}
	if dataType == 3 {
		defaultData = "字段名 | 类型 | 标签 (header) | 描述 | 是否必填 |\n"
	}
	// 表格标题
	md.WriteString(defaultData)
	md.WriteString("|--------|------|------------|------|------|\n")

	// 遍历字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		name := field.Name
		fieldType := field.Type.Name()
		Tag := field.Tag.Get("json")
		if dataType == 2 {
			Tag = field.Tag.Get("form")
		}
		if dataType == 3 {
			Tag = field.Tag.Get("header")
		}
		description := field.Tag.Get("desc")
		need := field.Tag.Get("required")
		require := ""
		if need == "true" {
			require = "是"
		}

		md.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", name, fieldType, Tag, description, require))

		if fieldValue.Kind() == reflect.Struct {
			md.WriteString(generateStructTable(fieldValue.Interface(), dataType, level+1))
		} else if fieldValue.Kind() == reflect.Slice && fieldValue.Len() > 0 {
			// 获取切片的第一个元素类型
			element := fieldValue.Index(0).Interface()
			md.WriteString(generateStructTable(element, dataType, level+1))
		}

	}

	return md.String()
}

func (this *CoreDb) GetData() map[string]APIGroup {
	var list []API
	this.DB.Find(&list)

	rs := make(map[string]APIGroup)

	for _, item := range list {
		// 检查分组是否存在
		group, ok := rs[item.Group]
		if !ok {
			// 如果分组不存在，初始化
			group = APIGroup{
				Group: item.Group,
				APIs:  make([]API, 0),
			}
		}

		// 修改分组中的 API 列表
		group.APIs = append(group.APIs, item)

		// 将修改后的分组放回 map
		rs[item.Group] = group
	}

	return rs

}
