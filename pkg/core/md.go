package core

const mdTemplate = `# 接口文档

{{- range . }}
## 接口分组：{{ .Group }}

{{- range .APIs }}
### 接口名称：{{ .Name }}
#### 请求路径：{{ .Path }}
#### 请求方式：{{ .Method }}
#### 请求参数：
{{ .Input }}
#### 请求结果：
{{ .Output }}
{{ end }}
{{ end }}
`

// 接口分组结构
type APIGroup struct {
	Group string
	APIs  []API
}
