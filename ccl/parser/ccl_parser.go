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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 35, 127,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 3, 2, 3, 2,
	3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 5, 3, 52, 10, 3, 3, 4, 3, 4, 3, 5, 3, 5,
	3, 5, 3, 5, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 5, 7, 65, 10, 7, 3, 8, 3, 8,
	3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 5, 8, 74, 10, 8, 3, 9, 3, 9, 3, 9, 3, 10,
	3, 10, 3, 10, 3, 11, 3, 11, 5, 11, 84, 10, 11, 3, 12, 3, 12, 3, 12, 3,
	12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 5, 13, 96, 10, 13, 3, 14,
	3, 14, 3, 15, 3, 15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3,
	18, 3, 18, 3, 19, 3, 19, 3, 20, 3, 20, 3, 20, 3, 20, 5, 20, 117, 10, 20,
	6, 20, 119, 10, 20, 13, 20, 14, 20, 120, 3, 21, 3, 21, 3, 22, 3, 22, 3,
	22, 2, 2, 23, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32,
	34, 36, 38, 40, 42, 2, 7, 3, 2, 25, 28, 3, 2, 10, 12, 3, 2, 14, 15, 4,
	2, 29, 29, 33, 34, 3, 2, 18, 24, 2, 115, 2, 44, 3, 2, 2, 2, 4, 51, 3, 2,
	2, 2, 6, 53, 3, 2, 2, 2, 8, 55, 3, 2, 2, 2, 10, 59, 3, 2, 2, 2, 12, 64,
	3, 2, 2, 2, 14, 73, 3, 2, 2, 2, 16, 75, 3, 2, 2, 2, 18, 78, 3, 2, 2, 2,
	20, 83, 3, 2, 2, 2, 22, 85, 3, 2, 2, 2, 24, 89, 3, 2, 2, 2, 26, 97, 3,
	2, 2, 2, 28, 99, 3, 2, 2, 2, 30, 101, 3, 2, 2, 2, 32, 103, 3, 2, 2, 2,
	34, 105, 3, 2, 2, 2, 36, 110, 3, 2, 2, 2, 38, 112, 3, 2, 2, 2, 40, 122,
	3, 2, 2, 2, 42, 124, 3, 2, 2, 2, 44, 45, 5, 4, 3, 2, 45, 46, 7, 3, 2, 2,
	46, 47, 5, 12, 7, 2, 47, 48, 7, 2, 2, 3, 48, 3, 3, 2, 2, 2, 49, 52, 5,
	6, 4, 2, 50, 52, 5, 8, 5, 2, 51, 49, 3, 2, 2, 2, 51, 50, 3, 2, 2, 2, 52,
	5, 3, 2, 2, 2, 53, 54, 5, 10, 6, 2, 54, 7, 3, 2, 2, 2, 55, 56, 5, 10, 6,
	2, 56, 57, 7, 4, 2, 2, 57, 58, 5, 12, 7, 2, 58, 9, 3, 2, 2, 2, 59, 60,
	7, 32, 2, 2, 60, 11, 3, 2, 2, 2, 61, 65, 5, 14, 8, 2, 62, 65, 5, 16, 9,
	2, 63, 65, 5, 34, 18, 2, 64, 61, 3, 2, 2, 2, 64, 62, 3, 2, 2, 2, 64, 63,
	3, 2, 2, 2, 65, 13, 3, 2, 2, 2, 66, 74, 5, 18, 10, 2, 67, 74, 5, 38, 20,
	2, 68, 74, 5, 20, 11, 2, 69, 70, 7, 5, 2, 2, 70, 71, 5, 12, 7, 2, 71, 72,
	7, 6, 2, 2, 72, 74, 3, 2, 2, 2, 73, 66, 3, 2, 2, 2, 73, 67, 3, 2, 2, 2,
	73, 68, 3, 2, 2, 2, 73, 69, 3, 2, 2, 2, 74, 15, 3, 2, 2, 2, 75, 76, 7,
	7, 2, 2, 76, 77, 5, 12, 7, 2, 77, 17, 3, 2, 2, 2, 78, 79, 7, 8, 2, 2, 79,
	80, 5, 10, 6, 2, 80, 19, 3, 2, 2, 2, 81, 84, 5, 22, 12, 2, 82, 84, 5, 24,
	13, 2, 83, 81, 3, 2, 2, 2, 83, 82, 3, 2, 2, 2, 84, 21, 3, 2, 2, 2, 85,
	86, 5, 10, 6, 2, 86, 87, 5, 42, 22, 2, 87, 88, 5, 40, 21, 2, 88, 23, 3,
	2, 2, 2, 89, 90, 5, 10, 6, 2, 90, 95, 5, 26, 14, 2, 91, 92, 5, 30, 16,
	2, 92, 93, 5, 32, 17, 2, 93, 96, 3, 2, 2, 2, 94, 96, 5, 28, 15, 2, 95,
	91, 3, 2, 2, 2, 95, 94, 3, 2, 2, 2, 96, 25, 3, 2, 2, 2, 97, 98, 9, 2, 2,
	2, 98, 27, 3, 2, 2, 2, 99, 100, 7, 9, 2, 2, 100, 29, 3, 2, 2, 2, 101, 102,
	7, 33, 2, 2, 102, 31, 3, 2, 2, 2, 103, 104, 9, 3, 2, 2, 104, 33, 3, 2,
	2, 2, 105, 106, 5, 14, 8, 2, 106, 107, 7, 13, 2, 2, 107, 108, 5, 36, 19,
	2, 108, 109, 5, 10, 6, 2, 109, 35, 3, 2, 2, 2, 110, 111, 9, 4, 2, 2, 111,
	37, 3, 2, 2, 2, 112, 113, 5, 10, 6, 2, 113, 118, 7, 16, 2, 2, 114, 116,
	5, 40, 21, 2, 115, 117, 7, 17, 2, 2, 116, 115, 3, 2, 2, 2, 116, 117, 3,
	2, 2, 2, 117, 119, 3, 2, 2, 2, 118, 114, 3, 2, 2, 2, 119, 120, 3, 2, 2,
	2, 120, 118, 3, 2, 2, 2, 120, 121, 3, 2, 2, 2, 121, 39, 3, 2, 2, 2, 122,
	123, 9, 5, 2, 2, 123, 41, 3, 2, 2, 2, 124, 125, 9, 6, 2, 2, 125, 43, 3,
	2, 2, 2, 9, 51, 64, 73, 83, 95, 116, 120,
}
var literalNames = []string{
	"", "'has'", "'with'", "'('", "')'", "'not'", "'empty'", "'now'", "'seconds'",
	"'days'", "'months'", "'in'", "'any'", "'all'", "'within'", "','", "'=='",
	"'!='", "'<='", "'<'", "'>'", "'>='", "'contains'", "'before'", "'after'",
	"'younger'", "'older'", "", "'true'", "'false'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "EqualsOperator",
	"NotEqualsOperator", "LessOrEqualsThanOperator", "LessThanOperator", "MoreThanOperator",
	"MoreOrEqualsThanOperator", "ContainsOperator", "BeforeOperator", "AfterOperator",
	"YoungerOperator", "OlderOperator", "BooleanLiteral", "True", "False",
	"Identifier", "Number", "StringLiteral", "Whitespace",
}

var ruleNames = []string{
	"condition", "assetType", "simpleAssetType", "filteredAssetType", "field",
	"expression", "simpleExpression", "notExpression", "emptyExpression", "comparison",
	"binaryComparison", "timeComparison", "timeOperator", "nowOperator", "time",
	"unit", "inExpression", "scope", "withinExpression", "value", "operator",
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
	CCLParserT__14                    = 15
	CCLParserEqualsOperator           = 16
	CCLParserNotEqualsOperator        = 17
	CCLParserLessOrEqualsThanOperator = 18
	CCLParserLessThanOperator         = 19
	CCLParserMoreThanOperator         = 20
	CCLParserMoreOrEqualsThanOperator = 21
	CCLParserContainsOperator         = 22
	CCLParserBeforeOperator           = 23
	CCLParserAfterOperator            = 24
	CCLParserYoungerOperator          = 25
	CCLParserOlderOperator            = 26
	CCLParserBooleanLiteral           = 27
	CCLParserTrue                     = 28
	CCLParserFalse                    = 29
	CCLParserIdentifier               = 30
	CCLParserNumber                   = 31
	CCLParserStringLiteral            = 32
	CCLParserWhitespace               = 33
)

// CCLParser rules.
const (
	CCLParserRULE_condition         = 0
	CCLParserRULE_assetType         = 1
	CCLParserRULE_simpleAssetType   = 2
	CCLParserRULE_filteredAssetType = 3
	CCLParserRULE_field             = 4
	CCLParserRULE_expression        = 5
	CCLParserRULE_simpleExpression  = 6
	CCLParserRULE_notExpression     = 7
	CCLParserRULE_emptyExpression   = 8
	CCLParserRULE_comparison        = 9
	CCLParserRULE_binaryComparison  = 10
	CCLParserRULE_timeComparison    = 11
	CCLParserRULE_timeOperator      = 12
	CCLParserRULE_nowOperator       = 13
	CCLParserRULE_time              = 14
	CCLParserRULE_unit              = 15
	CCLParserRULE_inExpression      = 16
	CCLParserRULE_scope             = 17
	CCLParserRULE_withinExpression  = 18
	CCLParserRULE_value             = 19
	CCLParserRULE_operator          = 20
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

func (s *ConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterCondition(s)
	}
}

func (s *ConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitCondition(s)
	}
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
		p.SetState(42)
		p.AssetType()
	}
	{
		p.SetState(43)
		p.Match(CCLParserT__0)
	}
	{
		p.SetState(44)
		p.Expression()
	}
	{
		p.SetState(45)
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

func (s *AssetTypeContext) FilteredAssetType() IFilteredAssetTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFilteredAssetTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFilteredAssetTypeContext)
}

func (s *AssetTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssetTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssetTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterAssetType(s)
	}
}

func (s *AssetTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitAssetType(s)
	}
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

	p.SetState(49)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(47)
			p.SimpleAssetType()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(48)
			p.FilteredAssetType()
		}

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

func (s *SimpleAssetTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterSimpleAssetType(s)
	}
}

func (s *SimpleAssetTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitSimpleAssetType(s)
	}
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
		p.SetState(51)
		p.Field()
	}

	return localctx
}

// IFilteredAssetTypeContext is an interface to support dynamic dispatch.
type IFilteredAssetTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFilteredAssetTypeContext differentiates from other interfaces.
	IsFilteredAssetTypeContext()
}

type FilteredAssetTypeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilteredAssetTypeContext() *FilteredAssetTypeContext {
	var p = new(FilteredAssetTypeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_filteredAssetType
	return p
}

func (*FilteredAssetTypeContext) IsFilteredAssetTypeContext() {}

func NewFilteredAssetTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilteredAssetTypeContext {
	var p = new(FilteredAssetTypeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_filteredAssetType

	return p
}

func (s *FilteredAssetTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *FilteredAssetTypeContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *FilteredAssetTypeContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FilteredAssetTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilteredAssetTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilteredAssetTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterFilteredAssetType(s)
	}
}

func (s *FilteredAssetTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitFilteredAssetType(s)
	}
}

func (p *CCLParser) FilteredAssetType() (localctx IFilteredAssetTypeContext) {
	localctx = NewFilteredAssetTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, CCLParserRULE_filteredAssetType)

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
		p.SetState(53)
		p.Field()
	}
	{
		p.SetState(54)
		p.Match(CCLParserT__1)
	}
	{
		p.SetState(55)
		p.Expression()
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

func (s *FieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterField(s)
	}
}

func (s *FieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitField(s)
	}
}

func (p *CCLParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, CCLParserRULE_field)

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
		p.SetState(57)
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

/*func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_expression
	return p
}*/

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

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *CCLParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, CCLParserRULE_expression)

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

	p.SetState(62)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(59)
			p.SimpleExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(60)
			p.NotExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(61)
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

func (s *SimpleExpressionContext) EmptyExpression() IEmptyExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEmptyExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEmptyExpressionContext)
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

func (s *SimpleExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterSimpleExpression(s)
	}
}

func (s *SimpleExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitSimpleExpression(s)
	}
}

func (p *CCLParser) SimpleExpression() (localctx ISimpleExpressionContext) {
	localctx = NewSimpleExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, CCLParserRULE_simpleExpression)

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

	p.SetState(71)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(64)
			p.EmptyExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(65)
			p.WithinExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(66)
			p.Comparison()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(67)
			p.Match(CCLParserT__2)
		}
		{
			p.SetState(68)
			p.Expression()
		}
		{
			p.SetState(69)
			p.Match(CCLParserT__3)
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

func (s *NotExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterNotExpression(s)
	}
}

func (s *NotExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitNotExpression(s)
	}
}

func (p *CCLParser) NotExpression() (localctx INotExpressionContext) {
	localctx = NewNotExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, CCLParserRULE_notExpression)

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
		p.SetState(73)
		p.Match(CCLParserT__4)
	}
	{
		p.SetState(74)
		p.Expression()
	}

	return localctx
}

// IEmptyExpressionContext is an interface to support dynamic dispatch.
type IEmptyExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEmptyExpressionContext differentiates from other interfaces.
	IsEmptyExpressionContext()
}

type EmptyExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEmptyExpressionContext() *EmptyExpressionContext {
	var p = new(EmptyExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_emptyExpression
	return p
}

func (*EmptyExpressionContext) IsEmptyExpressionContext() {}

func NewEmptyExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EmptyExpressionContext {
	var p = new(EmptyExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_emptyExpression

	return p
}

func (s *EmptyExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *EmptyExpressionContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *EmptyExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EmptyExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EmptyExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterEmptyExpression(s)
	}
}

func (s *EmptyExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitEmptyExpression(s)
	}
}

func (p *CCLParser) EmptyExpression() (localctx IEmptyExpressionContext) {
	localctx = NewEmptyExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, CCLParserRULE_emptyExpression)

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
		p.SetState(76)
		p.Match(CCLParserT__5)
	}
	{
		p.SetState(77)
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

func (s *ComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterComparison(s)
	}
}

func (s *ComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitComparison(s)
	}
}

func (p *CCLParser) Comparison() (localctx IComparisonContext) {
	localctx = NewComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, CCLParserRULE_comparison)

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

	p.SetState(81)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(79)
			p.BinaryComparison()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(80)
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

func (s *BinaryComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterBinaryComparison(s)
	}
}

func (s *BinaryComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitBinaryComparison(s)
	}
}

func (p *CCLParser) BinaryComparison() (localctx IBinaryComparisonContext) {
	localctx = NewBinaryComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, CCLParserRULE_binaryComparison)

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
		p.SetState(83)
		p.Field()
	}
	{
		p.SetState(84)
		p.Operator()
	}
	{
		p.SetState(85)
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

func (s *TimeComparisonContext) Time() ITimeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITimeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITimeContext)
}

func (s *TimeComparisonContext) Unit() IUnitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUnitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUnitContext)
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

func (s *TimeComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterTimeComparison(s)
	}
}

func (s *TimeComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitTimeComparison(s)
	}
}

func (p *CCLParser) TimeComparison() (localctx ITimeComparisonContext) {
	localctx = NewTimeComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, CCLParserRULE_timeComparison)

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
		p.Field()
	}
	{
		p.SetState(88)
		p.TimeOperator()
	}
	p.SetState(93)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CCLParserNumber:
		{
			p.SetState(89)
			p.Time()
		}
		{
			p.SetState(90)
			p.Unit()
		}

	case CCLParserT__6:
		{
			p.SetState(92)
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

func (s *TimeOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterTimeOperator(s)
	}
}

func (s *TimeOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitTimeOperator(s)
	}
}

func (p *CCLParser) TimeOperator() (localctx ITimeOperatorContext) {
	localctx = NewTimeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, CCLParserRULE_timeOperator)
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
		p.SetState(95)
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

func (s *NowOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterNowOperator(s)
	}
}

func (s *NowOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitNowOperator(s)
	}
}

func (p *CCLParser) NowOperator() (localctx INowOperatorContext) {
	localctx = NewNowOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, CCLParserRULE_nowOperator)

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
		p.SetState(97)
		p.Match(CCLParserT__6)
	}

	return localctx
}

// ITimeContext is an interface to support dynamic dispatch.
type ITimeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeContext differentiates from other interfaces.
	IsTimeContext()
}

type TimeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeContext() *TimeContext {
	var p = new(TimeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CCLParserRULE_time
	return p
}

func (*TimeContext) IsTimeContext() {}

func NewTimeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeContext {
	var p = new(TimeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CCLParserRULE_time

	return p
}

func (s *TimeContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeContext) Number() antlr.TerminalNode {
	return s.GetToken(CCLParserNumber, 0)
}

func (s *TimeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterTime(s)
	}
}

func (s *TimeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitTime(s)
	}
}

func (p *CCLParser) Time() (localctx ITimeContext) {
	localctx = NewTimeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, CCLParserRULE_time)

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
		p.Match(CCLParserNumber)
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

func (s *UnitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterUnit(s)
	}
}

func (s *UnitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitUnit(s)
	}
}

func (p *CCLParser) Unit() (localctx IUnitContext) {
	localctx = NewUnitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, CCLParserRULE_unit)
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
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CCLParserT__7)|(1<<CCLParserT__8)|(1<<CCLParserT__9))) != 0) {
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

func (s *InExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterInExpression(s)
	}
}

func (s *InExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitInExpression(s)
	}
}

func (p *CCLParser) InExpression() (localctx IInExpressionContext) {
	localctx = NewInExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, CCLParserRULE_inExpression)

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
		p.SetState(103)
		p.SimpleExpression()
	}
	{
		p.SetState(104)
		p.Match(CCLParserT__10)
	}
	{
		p.SetState(105)
		p.Scope()
	}
	{
		p.SetState(106)
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

func (s *ScopeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterScope(s)
	}
}

func (s *ScopeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitScope(s)
	}
}

func (p *CCLParser) Scope() (localctx IScopeContext) {
	localctx = NewScopeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, CCLParserRULE_scope)
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
		p.SetState(108)
		_la = p.GetTokenStream().LA(1)

		if !(_la == CCLParserT__11 || _la == CCLParserT__12) {
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

func (s *WithinExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterWithinExpression(s)
	}
}

func (s *WithinExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitWithinExpression(s)
	}
}

func (p *CCLParser) WithinExpression() (localctx IWithinExpressionContext) {
	localctx = NewWithinExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, CCLParserRULE_withinExpression)
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
		p.SetState(110)
		p.Field()
	}
	{
		p.SetState(111)
		p.Match(CCLParserT__13)
	}
	p.SetState(116)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-27)&-(0x1f+1)) == 0 && ((1<<uint((_la-27)))&((1<<(CCLParserBooleanLiteral-27))|(1<<(CCLParserNumber-27))|(1<<(CCLParserStringLiteral-27)))) != 0) {
		{
			p.SetState(112)
			p.Value()
		}
		p.SetState(114)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CCLParserT__14 {
			{
				p.SetState(113)
				p.Match(CCLParserT__14)
			}

		}

		p.SetState(118)
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

func (s *ValueContext) Number() antlr.TerminalNode {
	return s.GetToken(CCLParserNumber, 0)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterValue(s)
	}
}

func (s *ValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitValue(s)
	}
}

func (p *CCLParser) Value() (localctx IValueContext) {
	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, CCLParserRULE_value)
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
		p.SetState(120)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-27)&-(0x1f+1)) == 0 && ((1<<uint((_la-27)))&((1<<(CCLParserBooleanLiteral-27))|(1<<(CCLParserNumber-27))|(1<<(CCLParserStringLiteral-27)))) != 0) {
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

func (s *OperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.EnterOperator(s)
	}
}

func (s *OperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CCLListener); ok {
		listenerT.ExitOperator(s)
	}
}

func (p *CCLParser) Operator() (localctx IOperatorContext) {
	localctx = NewOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, CCLParserRULE_operator)
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
		p.SetState(122)
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
