package gdb

type DefaultHookModelInterfaceImpl struct {
	*Model
}

// TODO: 把所有接收者改为值类型实现
func (m *DefaultHookModelInterfaceImpl) setModel(model *Model) {
	m.Model = model
}

var (
	registerModelInterface = func(model *Model) ModelInterface {
		return &DefaultHookModelInterfaceImpl{
			Model: model,
		}
	}
)

func RegisterHookModelInterface(fn func(model *Model) ModelInterface) {
	if fn == nil {
		return
	}
	registerModelInterface = fn
}
