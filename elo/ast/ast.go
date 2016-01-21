// Abstract Syntax Tree

package ast

import (
)

type Node interface {
  Accept(v Visitor)
}

type Number struct {
  Value string
}

type Id struct {
  Value string
}

type Atom struct {
  Value string
}

type Keyword struct {
  Left  Node
  Right Node
}

type AtomKeyword struct {
  Left  Node
  Right Node
}

// TODO: atom positionals?
type CallArgs struct {
  Pos           []Node
  Keywords      []Node
  AtomKeywords  []Node
}

type Call struct {
  Left Node 
  Args Node 
}


func (node *Number) Accept(v Visitor) {
  v.VisitNumber(node)
}

func (node *Id) Accept(v Visitor) {
  v.VisitId(node)
}

func (node *Atom) Accept(v Visitor) {
  v.VisitAtom(node)
}

func (node *Keyword) Accept(v Visitor) {
  v.VisitKeyword(node)
}

func (node *AtomKeyword) Accept(v Visitor) {
  v.VisitAtomKeyword(node)
}

func (node *CallArgs) Accept(v Visitor) {
  v.VisitCallArgs(node)
}

func (node *Call) Accept(v Visitor) {
  v.VisitCall(node)
}