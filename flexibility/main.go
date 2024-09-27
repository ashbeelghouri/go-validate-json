package flexibility

import "github.com/ashbeelghouri/jsonschematics/utils"

type Types struct {
	TypeFunctions map[string]_Type
	Logger        utils.Logger
}

type _Type func(interface{}, map[string]interface{}) *interface{}

func (_t *Types) RegisterOperation(name string, fn _Type) {
	_t.Logger.DEBUG("registering types:", name)
	if _t.TypeFunctions == nil {
		_t.TypeFunctions = make(map[string]_Type)
	}
	_t.TypeFunctions[name] = fn
}

func (_t *Types) LoadBasicOperations() {
	_t.Logger.DEBUG("loading basic types")
	_t.RegisterOperation("IsString", IsString)
	_t.Logger.DEBUG("basic types loaded")
}
