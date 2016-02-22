package went

import (
)

const (
  MaxRegisters  = 249
  CallStackSize = 255
)

type callFrame struct {
  pc    int
  line  int
  fn    *Func
  r     [MaxRegisters]Value
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

func (state *State) RunProto(proto *FuncProto) {
  state.currentFrame = state.calls.New()
  state.currentFrame.fn = &Func{proto}

  execute(state)
}

func NewState() *State {
  return &State{
    Globals: make(map[string]Value, 128),
  }
}