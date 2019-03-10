package mg

import "path/filepath"

type mapContext struct {
	ctx map[string]interface{}
	str *string
}

// NewContext creates a simple [Context] based on a map.
func NewContext() Context {
	m := make(map[string]interface{}, 10)
	return &mapContext{ctx: m}
}

// ToContext converts the given map into a [Context].
//
// If given, the file is used in the String() representation of the Context.
func ToContext(m map[string]interface{}, file *ProcessedFile) Context {
	var str string
	if file != nil {
		path := file.Path
		if file.NewExtension != "" {
			path = changeFileExt(path, file.NewExtension)
		}
		// try to figure out a valid link to the file
		s, err := filepath.Rel("source/processed", path)
		if err == nil {
			str = "/" + s
		} else {
			str = path
		}
	}
	return &mapContext{ctx: m, str: &str}
}

var _ Context = (*mapContext)(nil)

func (m *mapContext) Get(name string) (interface{}, bool) {
	v, ok := m.ctx[name]
	return v, ok
}

func (m *mapContext) Remove(name string) {
	delete(m.ctx, name)
}

func (m *mapContext) Set(name string, value interface{}) {
	m.ctx[name] = value
}

func (m *mapContext) IsEmpty() bool {
	return len(m.ctx) == 0
}

func (m *mapContext) String() string {
	if m.str != nil {
		return *m.str
	}
	return m.String()
}