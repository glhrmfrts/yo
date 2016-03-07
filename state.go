package went

import (
	"github.com/glhrmfrts/went/parse"
)

const (
	MaxRegisters  = 249
	CallStackSize = 255
)

type callFrame struct {
	pc   int
	line int
	fn   *Func
	r    [MaxRegisters]Value
}

type callFrameStack struct {
	sp    int
	stack [CallStackSize]callFrame
}

type State struct {
	currentFrame *callFrame
	calls        callFrameStack
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
	state.currentFrame.fn = &Func{proto}

	execute(state)
	return Nil{}, nil
}

func NewState() *State {
	return &State{
		Globals: make(map[string]Value, 128),
	}
}
