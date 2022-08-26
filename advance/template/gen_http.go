package template

import (
	"io"
	"text/template"
)

type ServiceDefinition struct {
	Name    string
	Methods []Method
}

func (s *ServiceDefinition) GenName() string {
	return s.Name + "Gen"
}

type Method struct {
	Name         string
	ReqTypeName  string
	RespTypeName string
}

// 这是你们的作业，你们需要补全这个 template
const serviceTpl = `

`

func Gen(writer io.Writer, def *ServiceDefinition) error {
	tpl := template.New("service")
	tpl, err := tpl.Parse(serviceTpl)
	if err != nil {
		return err
	}
	// 还可以进一步调用 format.Source 来格式化生成代码
	return tpl.Execute(writer, def)
}
