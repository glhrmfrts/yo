package parse

// TODO: test files with lots of code

import (
	"fmt"
	"testing"
  //"github.com/glhrmfrts/elo-lang/elo/ast"
)

func TestExpr(t *testing.T) {
	valid := []string{
		"nil",
		"true",
		"false",
		"2",
		"2065",
		"3.141519",
		"8065.54e-6",
		"3e+5",
		"0.120213",
		".506",
		".5",
		"\"Hello world!\"",
		"\"Hello world\nI am a string!\"",
		"'Hello world\nI am a longer string with single quotes'",
		"identifier",
		"__identifier",
		"identIfier",
		"object.field",
		"array[2]",
		"object.field.with.array[2 + 5]",
		"array[5 + 9 * (12 / 24)].with.object.field",
		"calling()",
		"callingClosure()()",
		"calling().field",
		"object.field.calling()",
		"2 + 30 + 405",
		"2.1654 * 0.123 / 180e+1",
		"a*b-3/(5/2)",
		"(((((((5/2)))))))",
		"[1,2,3,4,5,'hello',6.45]",
		"[1,2,3,4,5,'hello',trailing,]",
		"{field: value,field2: func() {}}",
		"{2 + 2: value,\n\nfield: 'value\nvalue'}",
		"{field: 'value', trailing: true,}",
		"func() {}",
		"func(arg) { return arg }",
		"func(arg) { return 2 }",
		"func(arg) => arg * 2",
		"func(arg, arg2) => arg2, arg * 2 / 4",
		"func(a) ^(b) => a * b",
		"func(a) ^(b) ^(c) => a + b + c",
		"func(a) ^(b) { return a * b * 3 }",
		"func(a, b, c) ^(d, e, f, g) => a + b + c + d + e + f + g",
		"-2 + 5",
		"-(2 + 5)",
		"!false && true",
		"!(false && true)",
		"not false && true",
		"(not false) && true",
		"reallyLongNameWithALotOfUnsenseWordsOnIt",
		"5 < 2",
		"5 > 2",
		"a <= b * 2 / 3 * (4 ** 4)",
		"5 ** 5",
		"true ? 'is true' : 'is false'",
		"true ? 'is true' : true ? 'is still true' : 'is false'",
		"(98 < 100 ? 1 : 0) ? 'lt' : 'gt'",
	}

	fmt.Println("TestExpr:")
	for i, expr := range valid {
		_, err := ParseExpr([]byte(expr))
		if err != nil {
			t.Errorf("(%d) error: %s\n%s\n", i, expr, err.Error())
		} else {
			fmt.Printf("(%d) ok: %s\n\n", i, expr)
		}
	}
}