package templatedemo

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/likexian/gokit/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	type User struct {
		Name string
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, &User{
		Name: "John",
	})
	require.NoError(t, err)
	assert.Equal(t, "Hello, John", buffer.String())
}

func TestMapData(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, map[string]string{
		"Name": "John",
	})
	require.NoError(t, err)
	assert.Equal(t, "Hello, John", buffer.String())
}

func TestSlice(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{ index . 0}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, []string{
		"John",
	})
	require.NoError(t, err)
	assert.Equal(t, "Hello, John", buffer.String())
}

func TestBasic(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{.}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, "John")
	require.NoError(t, err)
	assert.Equal(t, "Hello, John", buffer.String())
}

func TestFuncCall(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
切片長度: {{len .Slice}}
{{printf "%.2f" 1.2345}}
Hello,{{.Hello "John" "Doe"}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, FuncCall{
		Slice: []string{"John", "Doe"},
	})
	require.NoError(t, err)
	assert.Equal(t, `
切片長度: 2
1.23
Hello, John.Doe`, buffer.String())
}

func TestLoop(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $index, $value := .Slice -}}
{{ . }}
{{$index}}-{{$value}}
{{- end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, "a\n0-ab\n1-b", buffer.String())
}

func TestLoop2(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range .Slice -}}
{{ . }}
{{- end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, "ab", buffer.String())
}

func TestForLoop(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $ele := . -}}
{{ $idx }}
{{- end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, make([]int, 5))
	require.NoError(t, err)
	assert.Equal(t, "01234", buffer.String())
}

func TestIfElseBlock(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (lt .Age 6) -}}
{{ .Age }} 是小孩
{{- else if and (gt .Age 6) (lt .Age 18) -}}
{{ .Age }} 是青少年
{{- else -}}
{{ .Age }} 是成年人
{{- end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	tpl.Execute(buffer, map[string]int{
		"Age": 10,
	})
	require.NoError(t, err)
	assert.Equal(t, "10 是青少年", buffer.String())
}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(first string, last string) string {
	return fmt.Sprintf(" %s.%s", first, last)
}

func TestPipeline(t *testing.T) {
	testCases := []struct {
		name string

		tpl  string
		data any

		want string
	}{
		// 这些例子来自官方文档
		// https://pkg.go.dev/text/template#hdr-Pipelines
		{
			name: "string constant",
			tpl:  `{{"\"output\""}}`,
			want: `"output"`,
		},
		{
			name: "raw string constant",
			tpl:  "{{`\"output\"`}}",
			want: `"output"`,
		},
		{
			name: "function call",
			tpl:  `{{printf "%q" "output"}}`,
			want: `"output"`,
		},
		{
			name: "take argument from pipeline",
			tpl:  `{{"output" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "parenthesized argument",
			tpl:  `{{printf "%q" (print "out" "put")}}`,
			want: `"output"`,
		},
		{
			name: "elaborate call",
			// printf "%s%s" "out" "put"
			tpl:  `{{"put" | printf "%s%s" "out" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "longer chain",
			tpl:  `{{"output" | printf "%s" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "with action using dot",
			tpl:  `{{with "output"}}{{printf "%q" .}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "with action that creates and uses a variable",
			tpl:  `{{with $x := "output" | printf "%q"}}{{$x}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "with action that uses the variable in another action",
			tpl:  `{{with $x := "output"}}{{printf "%q" $x}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "pipeline with action that uses the variable in another action",
			tpl:  `{{with $x := "output"}}{{$x | printf "%q"}}{{end}}`,
			want: `"output"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tpl := template.New(tc.name)
			tpl, err := tpl.Parse(tc.tpl)
			if err != nil {
				t.Fatal(err)
			}
			bs := &bytes.Buffer{}
			err = tpl.Execute(bs, tc.data)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, bs.String())
		})
	}
}
