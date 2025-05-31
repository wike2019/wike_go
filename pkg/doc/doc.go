package doc

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func Input(query interface{}, body interface{}, header interface{}) string {
	str := ""
	str += "query参数\n"
	str += "---\n"
	str += generateStructTable(query, 2, 1, false, true)
	str += "---\n"
	str += "body参数\n"
	str += "---\n"
	str += generateStructTable(body, 1, 1, false, true)
	str += "---\n"
	str += "header参数\n"
	str += "---\n"
	str += generateStructTable(header, 3, 1, false, true)
	str += "---\n"
	return str
}
func Output(body interface{}) string {
	return generateStructTable(body, 1, 1, false, false) + "---\n"
}

// 生成表格的函数
func generateStructTableIn(visited map[string]bool, data interface{}, dataType int, level int, Anonymous bool, short bool) string {
	if data == nil {
		return ""
	}
	stack := make([]string, 0)

	// 限制递归的最大深度
	const maxDepth = 7
	if level > maxDepth {
		return ""
	}

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// 如果是指针，解引用
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if v.IsValid() && !v.IsNil() {
			v = v.Elem()
		}
	}

	// 如果不是结构体，返回空
	if t.Kind() != reflect.Struct {
		return ""
	}

	md := strings.Builder{}
	if !Anonymous {
		if level > 5 {
			level = 5
		}
		// 添加标题
		md.WriteString(fmt.Sprintf("%s## %s\n\n", strings.Repeat("#", level), t.Name()))
	}

	// 根据 dataType 确定标签名称
	defaultData := "| 字段名 | 类型 | 标签 (json) | 描述 | 是否必填 | 是否搜索 |\n"
	switch dataType {
	case 2:
		defaultData = "| 字段名 | 类型 | 标签 (form) | 描述 | 是否必填 | 是否搜索 |\n"
	case 3:
		defaultData = "| 字段名 | 类型 | 标签 (header) | 描述 | 是否必填 | 是否搜索 |\n"
	}
	defaultDataLine := "|--------|------|------------|------|------|------|\n"
	if short == false {
		defaultData = "| 字段名 | 类型 | 标签 (json) | 描述  | \n"
		defaultDataLine = "|--------|------|------------|------| \n"
	}
	if !Anonymous {
		// 表格标题
		md.WriteString(defaultData)
		md.WriteString(defaultDataLine)
	}

	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 如果是匿名字段，将其字段展开到当前结构体
		if field.Anonymous {
			if fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Array || fieldValue.Kind() == reflect.Slice {
				// 使用结构体类型的名称来追踪已经访问过的结构体
				typeName := field.Type.String()
				if visited[typeName] && field.Tag.Get("deep") != "true" {
					// 如果该结构体已经访问过，则跳过
					continue
				}
				visited[typeName] = true
			}
			// 如果匿名字段是指针，需要解引用
			if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
				stack = append(stack, generateStructTableIn(visited, fieldValue.Elem().Interface(), dataType, level+1, field.Anonymous, short))
			} else if fieldValue.Kind() == reflect.Struct {
				stack = append(stack, generateStructTableIn(visited, fieldValue.Interface(), dataType, level+1, field.Anonymous, short))
			}
			continue
		}

		// 获取字段信息
		name := field.Name
		fieldType := field.Type.String()
		tag := field.Tag.Get("json")
		if dataType == 2 {
			tag = field.Tag.Get("form")
		} else if dataType == 3 {
			tag = field.Tag.Get("header")
		}
		description := field.Tag.Get("desc")
		if description == "" {
			description = field.Tag.Get("comment") // 支持 `comment` 标签
		}
		required := field.Tag.Get("required")
		require := "否"
		if required == "true" {
			require = "是"
		}
		search := "-"
		searched := field.Tag.Get("search")
		if searched != "" {
			search = "是"
		}

		// 判断是否是结构体或者数组类型

		if short == false {
			md.WriteString(fmt.Sprintf("| %s | %s | %s | %s  | \n", name, fieldType, tag, description))
		} else {
			// 添加字段到表格
			md.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | \n", name, fieldType, tag, description, require, search))
		}

		// 检查是否是标准库类型，如果是则跳过递归
		if isStandardLibraryType(field.Type) {

			continue
		}
		if fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Array || fieldValue.Kind() == reflect.Slice {
			// 使用结构体类型的名称来追踪已经访问过的结构体
			typeName := field.Type.String()
			if visited[typeName] && field.Tag.Get("deep") != "true" {
				// 如果该结构体已经访问过，则跳过
				continue
			}
			visited[typeName] = true
		}
		// 判断是否是可寻址的值类型

		//ptrAddr := fieldValue.Addr().Pointer()
		// 如果字段是指针类型，递归处理
		if fieldValue.Kind() == reflect.Ptr {
			// 如果是指针类型，解引用后递归处理
			if !fieldValue.IsNil() {
				stack = append(stack, generateStructTableIn(visited, fieldValue.Elem().Interface(), dataType, level+1, field.Anonymous, short))
			} else {
				if short == false {
					md.WriteString(fmt.Sprintf("| %s | %s | 空指针 | - ｜ \n", name, fieldType))
				} else {
					// 添加字段到表格
					md.WriteString(fmt.Sprintf("| %s | %s | - | 空指针 | - |\n", name, fieldType))
				}
			}
		} else if fieldValue.Kind() == reflect.Struct {
			// 如果是结构体类型，递归处理
			stack = append(stack, generateStructTableIn(visited, fieldValue.Interface(), dataType, level+1, field.Anonymous, short))
		} else if fieldValue.Kind() == reflect.Slice {
			// 如果是切片，递归处理第一个元素
			elementType := field.Type.Elem() // 获取切片的元素类型
			if fieldValue.Len() > 0 {
				element := fieldValue.Index(0).Interface()
				if short == false {
					md.WriteString(fmt.Sprintf("| %s | []%s | %s | %s ｜ \n", name, elementType.Name(), tag, description))
				} else {
					// 添加字段到表格
					md.WriteString(fmt.Sprintf("| %s | []%s | %s | %s | %s | %s | \n", name, elementType.Name(), tag, description, require, search))
				}
				stack = append(stack, generateStructTableIn(visited, element, dataType, level+1, field.Anonymous, short))
			} else {
				// 如果切片为空，递归生成切片元素类型
				emptyElement := reflect.New(elementType).Elem().Interface()
				stack = append(stack, generateStructTableIn(visited, emptyElement, dataType, level+1, field.Anonymous, short))
			}
			continue
		}

	}

	// 收集并返回最终结果
	res := md.String()
	for _, item := range stack {
		res += item
	}
	return res
}

// 生成表格的函数
func generateStructTable(data interface{}, dataType int, level int, Anonymous bool, short bool) string {
	if data == nil {
		return ""
	}
	stack := make([]string, 0)
	visited := make(map[string]bool) // 用于记录已访问过的结构体类型名称

	// 限制递归的最大深度
	const maxDepth = 7
	if level > maxDepth {
		return ""
	}

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// 如果是指针，解引用
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if v.IsValid() && !v.IsNil() {
			v = v.Elem()
		}
	}

	// 如果不是结构体，返回空
	if t.Kind() != reflect.Struct {
		return ""
	}

	md := strings.Builder{}
	if !Anonymous {
		if level > 5 {
			level = 5
		}
		// 添加标题
		md.WriteString(fmt.Sprintf("%s## %s\n\n", strings.Repeat("#", level), t.Name()))
	}

	// 根据 dataType 确定标签名称
	defaultData := "| 字段名 | 类型 | 标签 (json) | 描述 | 是否必填 | 是否搜索 |\n"
	switch dataType {
	case 2:
		defaultData = "| 字段名 | 类型 | 标签 (form) | 描述 | 是否必填 | 是否搜索 |\n"
	case 3:
		defaultData = "| 字段名 | 类型 | 标签 (header) | 描述 | 是否必填 | 是否搜索 |\n"
	}
	defaultDataLine := "|--------|------|------------|------|------|------|\n"
	if short == false {
		defaultData = "| 字段名 | 类型 | 标签 (json) | 描述  | \n"
		defaultDataLine = "|--------|------|------------|------| \n"
	}
	if !Anonymous {
		// 表格标题
		md.WriteString(defaultData)
		md.WriteString(defaultDataLine)
	}

	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 如果是匿名字段，将其字段展开到当前结构体
		if field.Anonymous {
			if fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Array || fieldValue.Kind() == reflect.Slice {
				// 使用结构体类型的名称来追踪已经访问过的结构体
				typeName := field.Type.String()
				if visited[typeName] && field.Tag.Get("deep") != "true" {
					// 如果该结构体已经访问过，则跳过
					continue
				}
				visited[typeName] = true
			}
			// 如果匿名字段是指针，需要解引用
			if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
				stack = append(stack, generateStructTableIn(visited, fieldValue.Elem().Interface(), dataType, level+1, field.Anonymous, short))
			} else if fieldValue.Kind() == reflect.Struct {
				stack = append(stack, generateStructTableIn(visited, fieldValue.Interface(), dataType, level+1, field.Anonymous, short))
			}
			continue
		}

		// 获取字段信息
		name := field.Name
		fieldType := field.Type.String()
		tag := field.Tag.Get("json")
		if dataType == 2 {
			tag = field.Tag.Get("form")
		} else if dataType == 3 {
			tag = field.Tag.Get("header")
		}
		description := field.Tag.Get("desc")
		if description == "" {
			description = field.Tag.Get("comment") // 支持 `comment` 标签
		}
		required := field.Tag.Get("required")
		require := "否"
		if required == "true" {
			require = "是"
		}
		search := "-"
		searched := field.Tag.Get("search")
		if searched != "" {
			search = "是"
		}

		// 判断是否是结构体或者数组类型

		if short == false {
			md.WriteString(fmt.Sprintf("| %s | %s | %s | %s  | \n", name, fieldType, tag, description))
		} else {
			// 添加字段到表格
			md.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | \n", name, fieldType, tag, description, require, search))
		}

		// 检查是否是标准库类型，如果是则跳过递归
		if isStandardLibraryType(field.Type) {

			continue
		}
		if fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Array || fieldValue.Kind() == reflect.Slice {
			// 使用结构体类型的名称来追踪已经访问过的结构体
			typeName := field.Type.String()
			if visited[typeName] && field.Tag.Get("deep") != "true" {
				// 如果该结构体已经访问过，则跳过
				continue
			}
			visited[typeName] = true
		}
		// 判断是否是可寻址的值类型

		//ptrAddr := fieldValue.Addr().Pointer()
		// 如果字段是指针类型，递归处理
		if fieldValue.Kind() == reflect.Ptr {
			// 如果是指针类型，解引用后递归处理
			if !fieldValue.IsNil() {
				stack = append(stack, generateStructTableIn(visited, fieldValue.Elem().Interface(), dataType, level+1, field.Anonymous, short))
			} else {
				if short == false {
					md.WriteString(fmt.Sprintf("| %s | %s | 空指针 | - ｜ \n", name, fieldType))
				} else {
					// 添加字段到表格
					md.WriteString(fmt.Sprintf("| %s | %s | - | 空指针 | - |\n", name, fieldType))
				}
			}
		} else if fieldValue.Kind() == reflect.Struct {
			// 如果是结构体类型，递归处理
			stack = append(stack, generateStructTableIn(visited, fieldValue.Interface(), dataType, level+1, field.Anonymous, short))
		} else if fieldValue.Kind() == reflect.Slice {
			// 如果是切片，递归处理第一个元素
			elementType := field.Type.Elem() // 获取切片的元素类型
			if fieldValue.Len() > 0 {
				element := fieldValue.Index(0).Interface()
				if short == false {
					md.WriteString(fmt.Sprintf("| %s | []%s | %s | %s ｜ \n", name, elementType.Name(), tag, description))
				} else {
					// 添加字段到表格
					md.WriteString(fmt.Sprintf("| %s | []%s | %s | %s | %s | %s | \n", name, elementType.Name(), tag, description, require, search))
				}
				stack = append(stack, generateStructTableIn(visited, element, dataType, level+1, field.Anonymous, short))
			} else {
				// 如果切片为空，递归生成切片元素类型
				emptyElement := reflect.New(elementType).Elem().Interface()
				stack = append(stack, generateStructTableIn(visited, emptyElement, dataType, level+1, field.Anonymous, short))
			}
			continue
		}

	}

	// 收集并返回最终结果
	res := md.String()
	for _, item := range stack {
		res += item
	}
	return res
}

func isStandardLibraryType(t reflect.Type) bool {
	// 如果类型是 nil，直接返回 false
	if t == nil {
		return false
	}
	if t == reflect.TypeOf(gorm.DeletedAt{}) {
		return true
	}
	// 对切片类型的元素类型检查
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	// 获取包路径
	pkgPath := t.PkgPath()

	// Go 标准库的类型通常在 "" 或 "time" 等包路径中
	if pkgPath == "" || pkgPath == "time" {
		return true
	}

	return false
}
