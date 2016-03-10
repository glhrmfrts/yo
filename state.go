package went

import (
	"fmt"
	"github.com/glhrmfrts/went/parse"
)

const (
	MaxRegisters  = 249
	CallStackSize = 255
)

type callFrame struct {
	pc         int
	line       int
	canRecover bool
	fn         *Func
	r          [MaxRegisters]Value
}

type callFrameStack struct {
	sp    int
	stack [CallStackSize]callFrame
}

type RuntimeError struct {
	Message string
	Line    int
	File    string
}

type State struct {
	currentFrame *callFrame
	calls        callFrameStack
	error        *RuntimeError
	Globals      map[string]Value
}

// callFrameStack

func (stack *callFrameStack) New() *callFrame {
	stack.sp += 1
	return &stack.stack[stack.sp-1]
}

func (stack *callFrameStack) Last() *callFrame {
	if stack.sp == 0 {
		return nil
	}
	return &stack.stack[stack.sp-1]
}

func (stack *callFrameStack) Rewind(n int) {
	for n > 0 {
		n--
		stack.sp--
	}
}

// callFrame

func (f *callFrame) recover() {
	// TODO: implement
}

// RuntimeError

func (err *RuntimeError) Error() string {
	return fmt.Sprintf("%s:%d: %s", err.File, err.Line, err.Message)
}

// State

func (state *State) DefineGlobal(name string, v Value) {
	state.Globals[name] = v
}

func (state *State) LoadString(source []byte, filename string) (Value, error) {
	nodes, err := parse.ParseFile(source, filename)
	if err != nil {
		return nil, err
	}

	proto, err := Compile(nodes, filename)
	if err != nil {
		return nil, err
	}

	return state.LoadProto(proto)
}

func (state *State) LoadProto(proto *FuncProto) (Value, error) {
	state.currentFrame = state.calls.New()
	state.currentFrame.fn = &Func{proto, false, nil}

	var err error
	if !execute(state) {
		err = state.error
	}

	return Nil{}, err
}

func (state *State) setError(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	state.error = &RuntimeError{Message: msg, Line: state.currentFrame.line, File: "test.elo"}
}

func NewState() *State {
	return &State{
		Globals: make(map[string]Value, 128),
	}
}
