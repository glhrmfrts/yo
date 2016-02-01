// Abstract Syntax Tree

package ast

import (
  "github.com/glhrmfrts/elo-lang/elo/token"
)

type Node interface {
  Accept(v Visitor)
}

type Nil struct {  
}

type Bool struct {
  Value bool
}

type Number struct {
  Value string
}

type Id struct {
  Value string
}

type String struct {
  Value string
}

type UnaryExpr struct {
  Op    token.Token
  Right Node
}

type BinaryExpr struct {
  Op    token.Token
  Left  Node
  Right Node
}


func (node *Nil) Accept(v Visitor) {
  v.VisitNil(node)
}

func (node *Bool) Accept(v Visitor) {
  v.VisitBool(node)
}

func (node *Number) Accept(v Visitor) {
  v.VisitNumber(node)
}

func (node *Id) Accept(v Visitor) {
  v.VisitId(node)
}

func (node *String) Accept(v Visitor) {
  v.VisitString(node)
}

func (node *UnaryExpr) Accept(v Visitor) {
  v.VisitUnaryExpr(node)
}

func (node *BinaryExpr) Accept(v Visitor) {
  v.VisitBinaryExpr(node)
}