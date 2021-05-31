// Code generated from CCL.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser // CCL

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 35, 118,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 3, 2, 3, 2, 3, 2, 3, 2, 3,
	2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 5, 6, 57, 10,
	6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 66, 10, 7, 3, 8, 3,
	8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 10, 3, 10, 5, 10, 76, 10, 10, 3, 11, 3, 11,
	3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12, 86, 10, 12, 3, 13, 3,
	13, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17,
	3, 17, 3, 17, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 108, 10,
	19, 6, 19, 110, 10, 19, 13, 19, 14, 19, 111, 3, 20, 3, 20, 3, 21, 3, 21,
	3, 21, 2, 2, 22, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
	32, 34, 36, 38, 40, 2, 7, 3, 2, 24, 27, 3, 2, 9, 11, 3, 2, 13, 14, 4, 2,
	28, 28, 32, 34, 3, 2, 17, 23, 2, 106, 2, 42, 3, 2, 2, 2, 4, 47, 3, 2, 2,
	2, 6, 49, 3, 2, 2, 2, 8, 51, 3, 2, 2, 2, 10, 56, 3, 2, 2, 2, 12, 65, 3,
	2, 2, 2, 14, 67, 3, 2, 2, 2, 16, 70, 3, 2, 2, 2, 18, 75, 3, 2, 2, 2, 20,
	77, 3, 2, 2, 2, 22, 81, 3, 2, 2, 2, 24, 87, 3, 2, 2, 2, 26, 89, 3, 2, 2,
	2, 28, 91, 3, 2, 2, 2, 30, 94, 3, 2, 2, 2, 32, 96, 3, 2, 2, 2, 34, 101,
	3, 2, 2, 2, 36, 103, 3, 2, 2, 2, 38, 113, 3, 2, 2, 2, 40, 115, 3, 2, 2,
	2, 42, 43, 5, 4, 3, 2, 43, 44, 7, 3, 2, 2, 44, 45, 5, 10, 6, 2, 45, 46,
	7, 2, 2, 3, 46, 3, 3, 2, 2, 2, 47, 48, 5, 6, 4, 2, 48, 5, 3, 2, 2, 2, 49,
	50, 5, 8, 5, 2, 50, 7, 3, 2, 2, 2, 51, 52, 7, 31, 2, 2, 52, 9, 3, 2, 2,
	2, 53, 57, 5, 12, 7, 2, 54, 57, 5, 14, 8, 2, 55, 57, 5, 32, 17, 2, 56,
	53, 3, 2, 2, 2, 56, 54, 3, 2, 2, 2, 56, 55, 3, 2, 2, 2, 57, 11, 3, 2, 2,
	2, 58, 66, 5, 16, 9, 2, 59, 66, 5, 36, 19, 2, 60, 66, 5, 18, 10, 2, 61,
	62, 7, 4, 2, 2, 62, 63, 5, 10, 6, 2, 63, 64, 7, 5, 2, 2, 64, 66, 3, 2,
	2, 2, 65, 58, 3, 2, 2, 2, 65, 59, 3, 2, 2, 2, 65, 60, 3, 2, 2, 2, 65, 61,
	3, 2, 2, 2, 66, 13, 3, 2, 2, 2, 67, 68, 7, 6, 2, 2, 68, 69, 5, 10, 6, 2,
	69, 15, 3, 2, 2, 2, 70, 71, 7, 7, 2, 2, 71, 72, 5, 8, 5, 2, 72, 17, 3,
	2, 2, 2, 73, 76, 5, 20, 11, 2, 74, 76, 5, 22, 12, 2, 75, 73, 3, 2, 2, 2,
	75, 74, 3, 2, 2, 2, 76, 19, 3, 2, 2, 2, 77, 78, 5, 8, 5, 2, 78, 79, 5,
	40, 21, 2, 79, 80, 5, 38, 20, 2, 80, 21, 3, 2, 2, 2, 81, 82, 5, 8, 5, 2,
	82, 85, 5, 24, 13, 2, 83, 86, 5, 28, 15, 2, 84, 86, 5, 26, 14, 2, 85, 83,
	3, 2, 2, 2, 85, 84, 3, 2, 2, 2, 86, 23, 3, 2, 2, 2, 87, 88, 9, 2, 2, 2,
	88, 25, 3, 2, 2, 2, 89, 90, 7, 8, 2, 2, 90, 27, 3, 2, 2, 2, 91, 92, 7,
	32, 2, 2, 92, 93, 5, 30, 16, 2, 93, 29, 3, 2, 2, 2, 94, 95, 9, 3, 2, 2,
	95, 31, 3, 2, 2, 2, 96, 97, 5, 12, 7, 2, 97, 98, 7, 12, 2, 2, 98, 99, 5,
	34, 18, 2, 99, 100, 5, 8, 5, 2, 100, 33, 3, 2, 2, 2, 101, 102, 9, 4, 2,
	2, 102, 35, 3, 2, 2, 2, 103, 104, 5, 8, 5, 2, 104, 109, 7, 15, 2, 2, 105,
	107, 5, 38, 20, 2, 106, 108, 7, 16, 2, 2, 107, 106, 3, 2, 2, 2, 107, 108,
	3, 2, 2, 2, 108, 110, 3, 2, 2, 2, 109, 105, 3, 2, 2, 2, 110, 111, 3, 2,
	2, 2, 111, 109, 3, 2, 2, 2, 111, 112, 3, 2, 2, 2, 112, 37, 3, 2, 2, 2,
	113, 114, 9, 5, 2, 2, 114, 39, 3, 2, 2, 2, 115, 116, 9, 6, 2, 2, 116, 41,
	3, 2, 2, 2, 8, 56, 65, 75, 85, 107, 111,
}
var literalNames = []string{
	"", "'has'", "'('", "')'", "'not'", "'empty'", "'now'", "'seconds'", "'days'",
	"'months'", "'in'", "'any'", "'all'", "'within'", "','", "'=='", "'!='",
	"'<='", "'<'", "'>'", "'>='", "'contains'", "'before'", "'after'", "'younger'",
	"'older'", "", "'true'", "'false'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "EqualsOperator",
	"NotEqualsOperator", "LessOrEqualsThanOperator", "LessThanOperator", "MoreThanOperator",
	"MoreOrEqualsThanOperator", "ContainsOperator", "BeforeOperator", "AfterOperator",
	"YoungerOperator", "OlderOperator", "BooleanLiteral", "True", "False",
	"Identifier", "IntNumber", "FloatNumber", "StringLiteral", "Whitespace",
}

var ruleNames = []string{
	"condition", "assetType", "simpleAssetType", "field", "expression", "simpleExpression",
	"notExpression", "isEmptyExpression", "comparison", "binaryComparison",
	"timeComparison", "timeOperator", "nowOperator", "duration", "unit", "inExpression",
	"scope", "withinExpression", "value", "operator",
}

type CCLParser struct {
	*antlr.BaseParser
}

// NewCCLParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *CCLParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewCCLParser(input antlr.TokenStream) *CCLParser {
	this := new(CCLParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "CCL.g4"

	return this
}

// CCLParser tokens.
const (
	CCLParserEOF                      = antlr.TokenEOF
	CCLParserT__0                     = 1
	CCLParserT__1                     = 2
	CCLParserT__2                     = 3
	CCLParserT__3                     = 4
	CCLParserT__4                     = 5
	CCLParserT__5                     = 6
	CCLParserT__6                     = 7
	CCLParserT__7                     = 8
	CCLParserT__8                     = 9
	CCLParserT__9                     = 10
	CCLParserT__10                    = 11
	CCLParserT__11                    = 12
	CCLParserT__12                    = 13
	CCLParserT__13                    = 14
	CCLParserEqualsOperator           = 15
	CCLParserNotEqualsOperator        = 16
	CCLParserLessOrEqualsThanOperator = 17
	CCLParserLessThanOperator         = 18
	CCLParserMoreThanOperator         = 19
	CCLParserMoreOrEqualsThanOperator = 20
	CCLParserContainsOperator         = 21
	CCLParserBeforeOperator           = 22
	CCLParserAfterOperator            = 23
	CCLParserYoungerOperator          = 24
	CCLParserOlderOperator            = 25
	CCLParserBooleanLiteral           = 26
	CCLParserTrue                     = 27
	CCLParserFalse                    = 28
	CCLParserIdentifier               = 29
	CCLParserIntNumber                = 30
	CCLParserFloatNumber              = 31
	CCLParserStringLiteral            = 32
	CCLParserWhitespace               = 33
)

// CCLParser rules.
const (
	CCLParserRULE_condition         = 0
	CCLParserRULE_assetType         = 1
	CCLParserRULE_simpleAssetType   = 2
	CCLParserRULE_field             = 3
	CCLParserRULE_expression        = 4
	CCLParserRULE_simpleExpression  = 5
	CCLParserRULE_notExpression     = 6
	CCLParserRULE_isEmptyExpression = 7
	CCLParserRULE_comparison        = 8
	CCLParserRULE_binaryComparison  = 9
	CCLParserRULE_timeComparison    = 10
	CCLParserRULE_timeOperator      = 11
	CCLParserRULE_nowOperator       = 12
	CCLParserRULE_duration          = 13
	CCLParserRULE_unit              = 14
	CCLParserRULE_inExpression      = 15
	CCLParserRULE_scope             = 16
	CCLParserRULE_withinExpression  = 17
	CCLParserRULE_value             = 18
	CCLParserRULE_operator          = 19
)

// IConditionContext is an interface to support dynamic dispatch.
type IConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsConditionContext differentiates from other interfaces.
	IsConditionContext()
}

type ConditionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConditionContext() *ConditionContext {
	var p = new(ConditionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_condition
	return p
}

func (*ConditionContext) IsConditionContext() {}

func NewConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionContext {
	var p = new(ConditionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_condition

	return p
}

func (s *ConditionContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionContext) AssetType() IAssetTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAssetTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAssetTypeContext)
}

func (s *ConditionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ConditionContext) EOF() antlr.TerminalNode {
	return s.GetToken(CCLParserEOF, 0)
}

func (s *ConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Condition() (localctx IConditionContext) {
	localctx = NewConditionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, CCLParserRULE_condition)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(40)
		p.AssetType()
	}
	{
		p.SetState(41)
		p.Match(CCLParserT__0)
	}
	{
		p.SetState(42)
		p.Expression()
	}
	{
		p.SetState(43)
		p.Match(CCLParserEOF)
	}

	return localctx
}

// IAssetTypeContext is an interface to support dynamic dispatch.
type IAssetTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAssetTypeContext differentiates from other interfaces.
	IsAssetTypeContext()
}

type AssetTypeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAssetTypeContext() *AssetTypeContext {
	var p = new(AssetTypeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_assetType
	return p
}

func (*AssetTypeContext) IsAssetTypeContext() {}

func NewAssetTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssetTypeContext {
	var p = new(AssetTypeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_assetType

	return p
}

func (s *AssetTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *AssetTypeContext) SimpleAssetType() ISimpleAssetTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISimpleAssetTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISimpleAssetTypeContext)
}

func (s *AssetTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssetTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) AssetType() (localctx IAssetTypeContext) {
	localctx = NewAssetTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, CCLParserRULE_assetType)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(45)
		p.SimpleAssetType()
	}

	return localctx
}

// ISimpleAssetTypeContext is an interface to support dynamic dispatch.
type ISimpleAssetTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSimpleAssetTypeContext differentiates from other interfaces.
	IsSimpleAssetTypeContext()
}

type SimpleAssetTypeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySimpleAssetTypeContext() *SimpleAssetTypeContext {
	var p = new(SimpleAssetTypeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_simpleAssetType
	return p
}

func (*SimpleAssetTypeContext) IsSimpleAssetTypeContext() {}

func NewSimpleAssetTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SimpleAssetTypeContext {
	var p = new(SimpleAssetTypeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_simpleAssetType

	return p
}

func (s *SimpleAssetTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *SimpleAssetTypeContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *SimpleAssetTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleAssetTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) SimpleAssetType() (localctx ISimpleAssetTypeContext) {
	localctx = NewSimpleAssetTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, CCLParserRULE_simpleAssetType)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(47)
		p.Field()
	}

	return localctx
}

// IFieldContext is an interface to support dynamic dispatch.
type IFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldContext differentiates from other interfaces.
	IsFieldContext()
}

type FieldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldContext() *FieldContext {
	var p = new(FieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_field
	return p
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	var p = new(FieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) Identifier() antlr.TerminalNode {
	return s.GetToken(CCLParserIdentifier, 0)
}

func (s *FieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, CCLParserRULE_field)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(49)
		p.Match(CCLParserIdentifier)
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) SimpleExpression() ISimpleExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISimpleExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISimpleExpressionContext)
}

func (s *ExpressionContext) NotExpression() INotExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INotExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INotExpressionContext)
}

func (s *ExpressionContext) InExpression() IInExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInExpressionContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, CCLParserRULE_expression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(54)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(51)
			p.SimpleExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(52)
			p.NotExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(53)
			p.InExpression()
		}

	}

	return localctx
}

// ISimpleExpressionContext is an interface to support dynamic dispatch.
type ISimpleExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSimpleExpressionContext differentiates from other interfaces.
	IsSimpleExpressionContext()
}

type SimpleExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySimpleExpressionContext() *SimpleExpressionContext {
	var p = new(SimpleExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_simpleExpression
	return p
}

func (*SimpleExpressionContext) IsSimpleExpressionContext() {}

func NewSimpleExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SimpleExpressionContext {
	var p = new(SimpleExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_simpleExpression

	return p
}

func (s *SimpleExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *SimpleExpressionContext) IsEmptyExpression() IIsEmptyExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIsEmptyExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIsEmptyExpressionContext)
}

func (s *SimpleExpressionContext) WithinExpression() IWithinExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithinExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithinExpressionContext)
}

func (s *SimpleExpressionContext) Comparison() IComparisonContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComparisonContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IComparisonContext)
}

func (s *SimpleExpressionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SimpleExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) SimpleExpression() (localctx ISimpleExpressionContext) {
	localctx = NewSimpleExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, CCLParserRULE_simpleExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(63)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(56)
			p.IsEmptyExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(57)
			p.WithinExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(58)
			p.Comparison()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(59)
			p.Match(CCLParserT__1)
		}
		{
			p.SetState(60)
			p.Expression()
		}
		{
			p.SetState(61)
			p.Match(CCLParserT__2)
		}

	}

	return localctx
}

// INotExpressionContext is an interface to support dynamic dispatch.
type INotExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNotExpressionContext differentiates from other interfaces.
	IsNotExpressionContext()
}

type NotExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNotExpressionContext() *NotExpressionContext {
	var p = new(NotExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_notExpression
	return p
}

func (*NotExpressionContext) IsNotExpressionContext() {}

func NewNotExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NotExpressionContext {
	var p = new(NotExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_notExpression

	return p
}

func (s *NotExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *NotExpressionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *NotExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NotExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) NotExpression() (localctx INotExpressionContext) {
	localctx = NewNotExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, CCLParserRULE_notExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(65)
		p.Match(CCLParserT__3)
	}
	{
		p.SetState(66)
		p.Expression()
	}

	return localctx
}

// IIsEmptyExpressionContext is an interface to support dynamic dispatch.
type IIsEmptyExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIsEmptyExpressionContext differentiates from other interfaces.
	IsIsEmptyExpressionContext()
}

type IsEmptyExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIsEmptyExpressionContext() *IsEmptyExpressionContext {
	var p = new(IsEmptyExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_isEmptyExpression
	return p
}

func (*IsEmptyExpressionContext) IsIsEmptyExpressionContext() {}

func NewIsEmptyExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IsEmptyExpressionContext {
	var p = new(IsEmptyExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_isEmptyExpression

	return p
}

func (s *IsEmptyExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *IsEmptyExpressionContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *IsEmptyExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IsEmptyExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) IsEmptyExpression() (localctx IIsEmptyExpressionContext) {
	localctx = NewIsEmptyExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, CCLParserRULE_isEmptyExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(68)
		p.Match(CCLParserT__4)
	}
	{
		p.SetState(69)
		p.Field()
	}

	return localctx
}

// IComparisonContext is an interface to support dynamic dispatch.
type IComparisonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsComparisonContext differentiates from other interfaces.
	IsComparisonContext()
}

type ComparisonContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonContext() *ComparisonContext {
	var p = new(ComparisonContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_comparison
	return p
}

func (*ComparisonContext) IsComparisonContext() {}

func NewComparisonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonContext {
	var p = new(ComparisonContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_comparison

	return p
}

func (s *ComparisonContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonContext) BinaryComparison() IBinaryComparisonContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBinaryComparisonContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBinaryComparisonContext)
}

func (s *ComparisonContext) TimeComparison() ITimeComparisonContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITimeComparisonContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITimeComparisonContext)
}

func (s *ComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Comparison() (localctx IComparisonContext) {
	localctx = NewComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, CCLParserRULE_comparison)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(73)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(71)
			p.BinaryComparison()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(72)
			p.TimeComparison()
		}

	}

	return localctx
}

// IBinaryComparisonContext is an interface to support dynamic dispatch.
type IBinaryComparisonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBinaryComparisonContext differentiates from other interfaces.
	IsBinaryComparisonContext()
}

type BinaryComparisonContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinaryComparisonContext() *BinaryComparisonContext {
	var p = new(BinaryComparisonContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_binaryComparison
	return p
}

func (*BinaryComparisonContext) IsBinaryComparisonContext() {}

func NewBinaryComparisonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryComparisonContext {
	var p = new(BinaryComparisonContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_binaryComparison

	return p
}

func (s *BinaryComparisonContext) GetParser() antlr.Parser { return s.parser }

func (s *BinaryComparisonContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *BinaryComparisonContext) Operator() IOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOperatorContext)
}

func (s *BinaryComparisonContext) Value() IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *BinaryComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryComparisonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) BinaryComparison() (localctx IBinaryComparisonContext) {
	localctx = NewBinaryComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, CCLParserRULE_binaryComparison)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(75)
		p.Field()
	}
	{
		p.SetState(76)
		p.Operator()
	}
	{
		p.SetState(77)
		p.Value()
	}

	return localctx
}

// ITimeComparisonContext is an interface to support dynamic dispatch.
type ITimeComparisonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeComparisonContext differentiates from other interfaces.
	IsTimeComparisonContext()
}

type TimeComparisonContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeComparisonContext() *TimeComparisonContext {
	var p = new(TimeComparisonContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_timeComparison
	return p
}

func (*TimeComparisonContext) IsTimeComparisonContext() {}

func NewTimeComparisonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeComparisonContext {
	var p = new(TimeComparisonContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_timeComparison

	return p
}

func (s *TimeComparisonContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeComparisonContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *TimeComparisonContext) TimeOperator() ITimeOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITimeOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITimeOperatorContext)
}

func (s *TimeComparisonContext) Duration() IDurationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationContext)
}

func (s *TimeComparisonContext) NowOperator() INowOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INowOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INowOperatorContext)
}

func (s *TimeComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeComparisonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) TimeComparison() (localctx ITimeComparisonContext) {
	localctx = NewTimeComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, CCLParserRULE_timeComparison)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(79)
		p.Field()
	}
	{
		p.SetState(80)
		p.TimeOperator()
	}
	p.SetState(83)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CCLParserIntNumber:
		{
			p.SetState(81)
			p.Duration()
		}

	case CCLParserT__5:
		{
			p.SetState(82)
			p.NowOperator()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ITimeOperatorContext is an interface to support dynamic dispatch.
type ITimeOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeOperatorContext differentiates from other interfaces.
	IsTimeOperatorContext()
}

type TimeOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeOperatorContext() *TimeOperatorContext {
	var p = new(TimeOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_timeOperator
	return p
}

func (*TimeOperatorContext) IsTimeOperatorContext() {}

func NewTimeOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeOperatorContext {
	var p = new(TimeOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_timeOperator

	return p
}

func (s *TimeOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeOperatorContext) BeforeOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserBeforeOperator, 0)
}

func (s *TimeOperatorContext) AfterOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserAfterOperator, 0)
}

func (s *TimeOperatorContext) YoungerOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserYoungerOperator, 0)
}

func (s *TimeOperatorContext) OlderOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserOlderOperator, 0)
}

func (s *TimeOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) TimeOperator() (localctx ITimeOperatorContext) {
	localctx = NewTimeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, CCLParserRULE_timeOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CCLParserBeforeOperator)|(1<<CCLParserAfterOperator)|(1<<CCLParserYoungerOperator)|(1<<CCLParserOlderOperator))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// INowOperatorContext is an interface to support dynamic dispatch.
type INowOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNowOperatorContext differentiates from other interfaces.
	IsNowOperatorContext()
}

type NowOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNowOperatorContext() *NowOperatorContext {
	var p = new(NowOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_nowOperator
	return p
}

func (*NowOperatorContext) IsNowOperatorContext() {}

func NewNowOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NowOperatorContext {
	var p = new(NowOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_nowOperator

	return p
}

func (s *NowOperatorContext) GetParser() antlr.Parser { return s.parser }
func (s *NowOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NowOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) NowOperator() (localctx INowOperatorContext) {
	localctx = NewNowOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, CCLParserRULE_nowOperator)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(87)
		p.Match(CCLParserT__5)
	}

	return localctx
}

// IDurationContext is an interface to support dynamic dispatch.
type IDurationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDurationContext differentiates from other interfaces.
	IsDurationContext()
}

type DurationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDurationContext() *DurationContext {
	var p = new(DurationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_duration
	return p
}

func (*DurationContext) IsDurationContext() {}

func NewDurationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DurationContext {
	var p = new(DurationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_duration

	return p
}

func (s *DurationContext) GetParser() antlr.Parser { return s.parser }

func (s *DurationContext) IntNumber() antlr.TerminalNode {
	return s.GetToken(CCLParserIntNumber, 0)
}

func (s *DurationContext) Unit() IUnitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUnitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUnitContext)
}

func (s *DurationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DurationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Duration() (localctx IDurationContext) {
	localctx = NewDurationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, CCLParserRULE_duration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(89)
		p.Match(CCLParserIntNumber)
	}
	{
		p.SetState(90)
		p.Unit()
	}

	return localctx
}

// IUnitContext is an interface to support dynamic dispatch.
type IUnitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUnitContext differentiates from other interfaces.
	IsUnitContext()
}

type UnitContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnitContext() *UnitContext {
	var p = new(UnitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_unit
	return p
}

func (*UnitContext) IsUnitContext() {}

func NewUnitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnitContext {
	var p = new(UnitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_unit

	return p
}

func (s *UnitContext) GetParser() antlr.Parser { return s.parser }
func (s *UnitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Unit() (localctx IUnitContext) {
	localctx = NewUnitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, CCLParserRULE_unit)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(92)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CCLParserT__6)|(1<<CCLParserT__7)|(1<<CCLParserT__8))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IInExpressionContext is an interface to support dynamic dispatch.
type IInExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInExpressionContext differentiates from other interfaces.
	IsInExpressionContext()
}

type InExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInExpressionContext() *InExpressionContext {
	var p = new(InExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_inExpression
	return p
}

func (*InExpressionContext) IsInExpressionContext() {}

func NewInExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InExpressionContext {
	var p = new(InExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_inExpression

	return p
}

func (s *InExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *InExpressionContext) SimpleExpression() ISimpleExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISimpleExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISimpleExpressionContext)
}

func (s *InExpressionContext) Scope() IScopeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IScopeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IScopeContext)
}

func (s *InExpressionContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *InExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) InExpression() (localctx IInExpressionContext) {
	localctx = NewInExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, CCLParserRULE_inExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(94)
		p.SimpleExpression()
	}
	{
		p.SetState(95)
		p.Match(CCLParserT__9)
	}
	{
		p.SetState(96)
		p.Scope()
	}
	{
		p.SetState(97)
		p.Field()
	}

	return localctx
}

// IScopeContext is an interface to support dynamic dispatch.
type IScopeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsScopeContext differentiates from other interfaces.
	IsScopeContext()
}

type ScopeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyScopeContext() *ScopeContext {
	var p = new(ScopeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_scope
	return p
}

func (*ScopeContext) IsScopeContext() {}

func NewScopeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScopeContext {
	var p = new(ScopeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_scope

	return p
}

func (s *ScopeContext) GetParser() antlr.Parser { return s.parser }
func (s *ScopeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScopeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Scope() (localctx IScopeContext) {
	localctx = NewScopeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, CCLParserRULE_scope)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(99)
		_la = p.GetTokenStream().LA(1)

		if !(_la == CCLParserT__10 || _la == CCLParserT__11) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IWithinExpressionContext is an interface to support dynamic dispatch.
type IWithinExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWithinExpressionContext differentiates from other interfaces.
	IsWithinExpressionContext()
}

type WithinExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithinExpressionContext() *WithinExpressionContext {
	var p = new(WithinExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_withinExpression
	return p
}

func (*WithinExpressionContext) IsWithinExpressionContext() {}

func NewWithinExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithinExpressionContext {
	var p = new(WithinExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_withinExpression

	return p
}

func (s *WithinExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *WithinExpressionContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *WithinExpressionContext) AllValue() []IValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IValueContext)(nil)).Elem())
	var tst = make([]IValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IValueContext)
		}
	}

	return tst
}

func (s *WithinExpressionContext) Value(i int) IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *WithinExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithinExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) WithinExpression() (localctx IWithinExpressionContext) {
	localctx = NewWithinExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, CCLParserRULE_withinExpression)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(101)
		p.Field()
	}
	{
		p.SetState(102)
		p.Match(CCLParserT__12)
	}
	p.SetState(107)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-26)&-(0x1f+1)) == 0 && ((1<<uint((_la-26)))&((1<<(CCLParserBooleanLiteral-26))|(1<<(CCLParserIntNumber-26))|(1<<(CCLParserFloatNumber-26))|(1<<(CCLParserStringLiteral-26)))) != 0) {
		{
			p.SetState(103)
			p.Value()
		}
		p.SetState(105)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CCLParserT__13 {
			{
				p.SetState(104)
				p.Match(CCLParserT__13)
			}

		}

		p.SetState(109)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IValueContext is an interface to support dynamic dispatch.
type IValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsValueContext differentiates from other interfaces.
	IsValueContext()
}

type ValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueContext() *ValueContext {
	var p = new(ValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_value
	return p
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) StringLiteral() antlr.TerminalNode {
	return s.GetToken(CCLParserStringLiteral, 0)
}

func (s *ValueContext) BooleanLiteral() antlr.TerminalNode {
	return s.GetToken(CCLParserBooleanLiteral, 0)
}

func (s *ValueContext) IntNumber() antlr.TerminalNode {
	return s.GetToken(CCLParserIntNumber, 0)
}

func (s *ValueContext) FloatNumber() antlr.TerminalNode {
	return s.GetToken(CCLParserFloatNumber, 0)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Value() (localctx IValueContext) {
	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, CCLParserRULE_value)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(111)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-26)&-(0x1f+1)) == 0 && ((1<<uint((_la-26)))&((1<<(CCLParserBooleanLiteral-26))|(1<<(CCLParserIntNumber-26))|(1<<(CCLParserFloatNumber-26))|(1<<(CCLParserStringLiteral-26)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IOperatorContext is an interface to support dynamic dispatch.
type IOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOperatorContext differentiates from other interfaces.
	IsOperatorContext()
}

type OperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperatorContext() *OperatorContext {
	var p = new(OperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_operator
	return p
}

func (*OperatorContext) IsOperatorContext() {}

func NewOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperatorContext {
	var p = new(OperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_operator

	return p
}

func (s *OperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *OperatorContext) EqualsOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserEqualsOperator, 0)
}

func (s *OperatorContext) NotEqualsOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserNotEqualsOperator, 0)
}

func (s *OperatorContext) LessOrEqualsThanOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserLessOrEqualsThanOperator, 0)
}

func (s *OperatorContext) LessThanOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserLessThanOperator, 0)
}

func (s *OperatorContext) MoreThanOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserMoreThanOperator, 0)
}

func (s *OperatorContext) MoreOrEqualsThanOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserMoreOrEqualsThanOperator, 0)
}

func (s *OperatorContext) ContainsOperator() antlr.TerminalNode {
	return s.GetToken(CCLParserContainsOperator, 0)
}

func (s *OperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *CCLParser) Operator() (localctx IOperatorContext) {
	localctx = NewOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, CCLParserRULE_operator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(113)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CCLParserEqualsOperator)|(1<<CCLParserNotEqualsOperator)|(1<<CCLParserLessOrEqualsThanOperator)|(1<<CCLParserLessThanOperator)|(1<<CCLParserMoreThanOperator)|(1<<CCLParserMoreOrEqualsThanOperator)|(1<<CCLParserContainsOperator))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}
