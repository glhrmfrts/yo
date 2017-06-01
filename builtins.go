package yo

import (
	"fmt"
)

func defineBuiltins(vm *VM) {
	vm.Define("append", GoFunc(builtinAppend))
	vm.Define("isnumber", GoFunc(builtinIsNumber))
	vm.Define("len", GoFunc(builtinLen))
	vm.Define("println", GoFunc(builtinPrintln))
	vm.Define("type", GoFunc(builtinType))
}

func builtinAppend(call *FuncCall) {
	if call.NumArgs == uint(0) {
		return
	}

	ptr := call.Args[0]
	arr := ptr.(*Array)
	*arr = append(*arr, call.Args[1:]...)

	call.PushReturnValue(ptr)
}

func builtinIsNumber(call *FuncCall) {
	if call.NumArgs <= uint(0) {
		call.PushReturnValue(Bool(false))
	} else {
		call.PushReturnValue(Bool(call.Args[0].Type() == ValueNumber))
	}
}

func builtinLen(call *FuncCall) {
	if call.NumArgs == uint(0) {
		panic("len expects 1 argument")
	} else {
		n := len(*(call.Args[0].(*Array)))
		call.PushReturnValue(Number(n))
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
