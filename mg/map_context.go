package mg

type mapContext struct {
	ctx map[string]interface{}
}

func NewContext() Context {
	m := make(map[string]interface{}, 10)
	return &mapContext{ctx: m}
}

func ToContext(m map[string]interface{}) Context {
	return &mapContext{ctx: m}
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
