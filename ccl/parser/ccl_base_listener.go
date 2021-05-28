// Code generated from CCL.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser // CCL

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseCCLListener is a complete listener for a parse tree produced by CCLParser.
type BaseCCLListener struct{}

var _ CCLListener = &BaseCCLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseCCLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseCCLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseCCLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseCCLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterCondition is called when production condition is entered.
func (s *BaseCCLListener) EnterCondition(ctx *ConditionContext) {}

// ExitCondition is called when production condition is exited.
func (s *BaseCCLListener) ExitCondition(ctx *ConditionContext) {}

// EnterAssetType is called when production assetType is entered.
func (s *BaseCCLListener) EnterAssetType(ctx *AssetTypeContext) {}

// ExitAssetType is called when production assetType is exited.
func (s *BaseCCLListener) ExitAssetType(ctx *AssetTypeContext) {}

// EnterSimpleAssetType is called when production simpleAssetType is entered.
func (s *BaseCCLListener) EnterSimpleAssetType(ctx *SimpleAssetTypeContext) {}

// ExitSimpleAssetType is called when production simpleAssetType is exited.
func (s *BaseCCLListener) ExitSimpleAssetType(ctx *SimpleAssetTypeContext) {}

// EnterFilteredAssetType is called when production filteredAssetType is entered.
func (s *BaseCCLListener) EnterFilteredAssetType(ctx *FilteredAssetTypeContext) {}

// ExitFilteredAssetType is called when production filteredAssetType is exited.
func (s *BaseCCLListener) ExitFilteredAssetType(ctx *FilteredAssetTypeContext) {}

// EnterField is called when production field is entered.
func (s *BaseCCLListener) EnterField(ctx *FieldContext) {}

// ExitField is called when production field is exited.
func (s *BaseCCLListener) ExitField(ctx *FieldContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseCCLListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseCCLListener) ExitExpression(ctx *ExpressionContext) {}

// EnterSimpleExpression is called when production simpleExpression is entered.
func (s *BaseCCLListener) EnterSimpleExpression(ctx *SimpleExpressionContext) {}

// ExitSimpleExpression is called when production simpleExpression is exited.
func (s *BaseCCLListener) ExitSimpleExpression(ctx *SimpleExpressionContext) {}

// EnterNotExpression is called when production notExpression is entered.
func (s *BaseCCLListener) EnterNotExpression(ctx *NotExpressionContext) {}

// ExitNotExpression is called when production notExpression is exited.
func (s *BaseCCLListener) ExitNotExpression(ctx *NotExpressionContext) {}

// EnterEmptyExpression is called when production emptyExpression is entered.
func (s *BaseCCLListener) EnterEmptyExpression(ctx *EmptyExpressionContext) {}

// ExitEmptyExpression is called when production emptyExpression is exited.
func (s *BaseCCLListener) ExitEmptyExpression(ctx *EmptyExpressionContext) {}

// EnterComparison is called when production comparison is entered.
func (s *BaseCCLListener) EnterComparison(ctx *ComparisonContext) {}

// ExitComparison is called when production comparison is exited.
func (s *BaseCCLListener) ExitComparison(ctx *ComparisonContext) {}

// EnterBinaryComparison is called when production binaryComparison is entered.
func (s *BaseCCLListener) EnterBinaryComparison(ctx *BinaryComparisonContext) {}

// ExitBinaryComparison is called when production binaryComparison is exited.
func (s *BaseCCLListener) ExitBinaryComparison(ctx *BinaryComparisonContext) {}

// EnterTimeComparison is called when production timeComparison is entered.
func (s *BaseCCLListener) EnterTimeComparison(ctx *TimeComparisonContext) {}

// ExitTimeComparison is called when production timeComparison is exited.
func (s *BaseCCLListener) ExitTimeComparison(ctx *TimeComparisonContext) {}

// EnterTimeOperator is called when production timeOperator is entered.
func (s *BaseCCLListener) EnterTimeOperator(ctx *TimeOperatorContext) {}

// ExitTimeOperator is called when production timeOperator is exited.
func (s *BaseCCLListener) ExitTimeOperator(ctx *TimeOperatorContext) {}

// EnterNowOperator is called when production nowOperator is entered.
func (s *BaseCCLListener) EnterNowOperator(ctx *NowOperatorContext) {}

// ExitNowOperator is called when production nowOperator is exited.
func (s *BaseCCLListener) ExitNowOperator(ctx *NowOperatorContext) {}

// EnterTime is called when production time is entered.
func (s *BaseCCLListener) EnterTime(ctx *TimeContext) {}

// ExitTime is called when production time is exited.
func (s *BaseCCLListener) ExitTime(ctx *TimeContext) {}

// EnterUnit is called when production unit is entered.
func (s *BaseCCLListener) EnterUnit(ctx *UnitContext) {}

// ExitUnit is called when production unit is exited.
func (s *BaseCCLListener) ExitUnit(ctx *UnitContext) {}

// EnterInExpression is called when production inExpression is entered.
func (s *BaseCCLListener) EnterInExpression(ctx *InExpressionContext) {}

// ExitInExpression is called when production inExpression is exited.
func (s *BaseCCLListener) ExitInExpression(ctx *InExpressionContext) {}

// EnterScope is called when production scope is entered.
func (s *BaseCCLListener) EnterScope(ctx *ScopeContext) {}

// ExitScope is called when production scope is exited.
func (s *BaseCCLListener) ExitScope(ctx *ScopeContext) {}

// EnterWithinExpression is called when production withinExpression is entered.
func (s *BaseCCLListener) EnterWithinExpression(ctx *WithinExpressionContext) {}

// ExitWithinExpression is called when production withinExpression is exited.
func (s *BaseCCLListener) ExitWithinExpression(ctx *WithinExpressionContext) {}

// EnterValue is called when production value is entered.
func (s *BaseCCLListener) EnterValue(ctx *ValueContext) {}

// ExitValue is called when production value is exited.
func (s *BaseCCLListener) ExitValue(ctx *ValueContext) {}

// EnterOperator is called when production operator is entered.
func (s *BaseCCLListener) EnterOperator(ctx *OperatorContext) {}

// ExitOperator is called when production operator is exited.
func (s *BaseCCLListener) ExitOperator(ctx *OperatorContext) {}
