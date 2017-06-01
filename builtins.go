package yo

import (
	"fmt"
)

func defineBuiltins(vm *VM) {
	vm.Define("append", GoFunc(builtinAppend))
	vm.Define("isnumber", GoFunc(builtinIsNumber))
	vm.Define("println", GoFunc(builtinPrintln))
	vm.Define("type", GoFunc(builtinType))
}

func builtinAppend(call *FuncCall) {
	if call.NumArgs <= uint(0) {
		return
	}

	arr := call.Args[0]
	call.PushReturnValue(append(arr.(Array), call.Args[1:]...))
}

func builtinIsNumber(call *FuncCall) {
	if call.NumArgs <= uint(0) {
		call.PushReturnValue(Bool(false))
	} else {
		call.PushReturnValue(Bool(call.Args[0].Type() == ValueNumber))
	}
}

func builtinPrintln(call *FuncCall) {
	for i := uint(0); i < call.NumArgs; i++ {
		fmt.Printf("%v", call.Args[i])
	}

	fmt.Println()
}

func builtinType(call *FuncCall) {
	if call.NumArgs <= uint(0) {
		call.PushReturnValue(String("nil"))
	} else {
		call.PushReturnValue(String(call.Args[0].Type().String()))
	}
}
