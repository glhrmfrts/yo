package elo

import (
)

const (
  MaxRegisters  = 249
  CallStackSize = 255
)

type callFrame struct {
  pc    int
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
}


func (stack *callFrameStack) New() *callFrame {
  stack.sp += 1
  return &stack.stack[stack.sp-1]
}

func (stack *callFrameStack) Back() *callFrame {
  if stack.sp == 0 {
    return nil
  }
  return &stack.stack[stack.sp-1]
}

func (state *State) RunProto(proto *FuncProto) {
  state.currentFrame = state.calls.New()
  state.currentFrame.fn = &Func{proto}

  execute(state)
}