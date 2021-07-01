package gen

import (
	"embed"
	"io"
	"text/template"
)

//go:embed "template/*.tmpl"
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "template/*.tmpl"))

type Import struct {
	Path  string
	Alias string
}

type Arg struct {
	Name string
	Type string
}

type CGoMsgSend struct {
	Name        string
	Selector    string
	MsgSendFunc string
	Args        []Arg
	Return      string
}

func (m CGoMsgSend) HasReturn() bool {
	return m.Return != "void"
}

type CGoWrapperArg struct {
	Name     string
	Type     string
	ToCGoFmt string
}

type CGoWrapperReturn struct {
	Type       string
	FromCGoFmt string
}

type CGoWrapperFunc struct {
	Name    string
	Args    []CGoWrapperArg
	Returns []CGoWrapperReturn
}

func (f CGoWrapperFunc) HasReturn() bool {
	return len(f.Returns) > 0
}

type MethodReturn struct {
	Type string
}

type MethodDef struct {
	Name        string
	Args        []Arg
	Returns     []MethodReturn
	WrappedFunc string
}

func (md MethodDef) HasReturn() bool {
	return len(md.Returns) > 0
}

type ClassDef struct {
	Name            string
	ClassMethods    []MethodDef
	InstanceMethods []MethodDef
}

type GoPackage struct {
	Package         string
	Imports         []Import
	MsgSendWrappers []CGoMsgSend
	CGoWrapperFuncs []CGoWrapperFunc
	Classes         []ClassDef
}

func (p *GoPackage) Generate(w io.Writer) error {
	return templates.ExecuteTemplate(w, "main", p)
}
