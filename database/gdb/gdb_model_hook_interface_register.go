package gdb

type DefaultModelInterfaceImpl struct {
	*Model
}

func (m DefaultModelInterfaceImpl) setModel(model *Model) {
	m.Model = model
}

var (
	registerModelInterface = func(model *Model) ModelInterface {
		return DefaultModelInterfaceImpl{
			Model: model,
		}
	}
)

func RegisterModelInterface(fn func(model *Model) ModelInterface) {
	if fn == nil {
		return
	}
	registerModelInterface = fn
}
