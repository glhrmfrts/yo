package vm

type compiler struct {
  filename string
}

//
// visitor interface functions
//

func (c *compiler) VisitNil(node *ast.Nil) {

}

func (c *compiler) VisitBool(node *ast.Bool) {

}

func (c *compiler) VisitNumber(node *ast.Number) {

}

func (c *compiler) VisitId(node *ast.Id) {

}

func (c *compiler) VisitString(node *ast.String) {

}

func (c *compiler) VisitArray(node *ast.Array) {

}

func (c *compiler) VisitObjectField(node *ast.ObjectField) {

}

func (c *compiler) VisitObject(node *ast.Object) {

}

func (c *compiler) VisitFunction(node *ast.Function) {
 
}

func (c *compiler) VisitSelector(node *ast.Selector) {
 
}

func (c *compiler) VisitSubscript(node *ast.Subscript) {

}

func (c *compiler) VisitSlice(node *ast.Slice) {

}

func (c *compiler) VisitKwArg(node *ast.KwArg) {
  
}

func (c *compiler) VisitVarArg(node *ast.VarArg) {

}

func (c *compiler) VisitCallExpr(node *ast.CallExpr) {

}

func (c *compiler) VisitUnaryExpr(node *ast.UnaryExpr) {
 
}

func (c *compiler) VisitBinaryExpr(node *ast.BinaryExpr) {

}

func (c *compiler) VisitDeclaration(node *ast.Declaration) {
 
}

func (c *compiler) VisitAssignment(node *ast.Assignment) {
 
}

func (c *compiler) VisitBranchStmt(node *ast.BranchStmt) {

}

func (c *compiler) VisitReturnStmt(node *ast.ReturnStmt) {

}

func (c *compiler) VisitIfStmt(node *ast.IfStmt) {
 
}

func (c *compiler) VisitForIteratorStmt(node *ast.ForIteratorStmt) {

}

func (c *compiler) VisitForStmt(node *ast.ForStmt) {

}

func (c *compiler) VisitBlock(node *ast.Block) {

}

func Compile(root ast.Node, filename string) {
  var c compiler
  c.filename = filename
  root.Accept(c)
}