// ccl contains the the Cloud Compliance Language (CCL)
package ccl

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"clouditor.io/clouditor/ccl/parser"
)

func init() {

}

type TreeShapeListener struct {
	*parser.BaseCCLListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (this *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

func RunRule() {
	input, _ := antlr.NewFileStream("../rules/encryption.ccl")
	lexer := parser.NewCCLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewCCLParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Condition()
	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}
