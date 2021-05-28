// Code generated from CCL.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser // CCL

import "github.com/antlr/antlr4/runtime/Go/antlr"

// CCLListener is a complete listener for a parse tree produced by CCLParser.
type CCLListener interface {
	antlr.ParseTreeListener

	// EnterCondition is called when entering the condition production.
	EnterCondition(c *ConditionContext)

	// EnterAssetType is called when entering the assetType production.
	EnterAssetType(c *AssetTypeContext)

	// EnterSimpleAssetType is called when entering the simpleAssetType production.
	EnterSimpleAssetType(c *SimpleAssetTypeContext)

	// EnterFilteredAssetType is called when entering the filteredAssetType production.
	EnterFilteredAssetType(c *FilteredAssetTypeContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterSimpleExpression is called when entering the simpleExpression production.
	EnterSimpleExpression(c *SimpleExpressionContext)

	// EnterNotExpression is called when entering the notExpression production.
	EnterNotExpression(c *NotExpressionContext)

	// EnterEmptyExpression is called when entering the emptyExpression production.
	EnterEmptyExpression(c *EmptyExpressionContext)

	// EnterComparison is called when entering the comparison production.
	EnterComparison(c *ComparisonContext)

	// EnterBinaryComparison is called when entering the binaryComparison production.
	EnterBinaryComparison(c *BinaryComparisonContext)

	// EnterTimeComparison is called when entering the timeComparison production.
	EnterTimeComparison(c *TimeComparisonContext)

	// EnterTimeOperator is called when entering the timeOperator production.
	EnterTimeOperator(c *TimeOperatorContext)

	// EnterNowOperator is called when entering the nowOperator production.
	EnterNowOperator(c *NowOperatorContext)

	// EnterTime is called when entering the time production.
	EnterTime(c *TimeContext)

	// EnterUnit is called when entering the unit production.
	EnterUnit(c *UnitContext)

	// EnterInExpression is called when entering the inExpression production.
	EnterInExpression(c *InExpressionContext)

	// EnterScope is called when entering the scope production.
	EnterScope(c *ScopeContext)

	// EnterWithinExpression is called when entering the withinExpression production.
	EnterWithinExpression(c *WithinExpressionContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

	// EnterOperator is called when entering the operator production.
	EnterOperator(c *OperatorContext)

	// ExitCondition is called when exiting the condition production.
	ExitCondition(c *ConditionContext)

	// ExitAssetType is called when exiting the assetType production.
	ExitAssetType(c *AssetTypeContext)

	// ExitSimpleAssetType is called when exiting the simpleAssetType production.
	ExitSimpleAssetType(c *SimpleAssetTypeContext)

	// ExitFilteredAssetType is called when exiting the filteredAssetType production.
	ExitFilteredAssetType(c *FilteredAssetTypeContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitSimpleExpression is called when exiting the simpleExpression production.
	ExitSimpleExpression(c *SimpleExpressionContext)

	// ExitNotExpression is called when exiting the notExpression production.
	ExitNotExpression(c *NotExpressionContext)

	// ExitEmptyExpression is called when exiting the emptyExpression production.
	ExitEmptyExpression(c *EmptyExpressionContext)

	// ExitComparison is called when exiting the comparison production.
	ExitComparison(c *ComparisonContext)

	// ExitBinaryComparison is called when exiting the binaryComparison production.
	ExitBinaryComparison(c *BinaryComparisonContext)

	// ExitTimeComparison is called when exiting the timeComparison production.
	ExitTimeComparison(c *TimeComparisonContext)

	// ExitTimeOperator is called when exiting the timeOperator production.
	ExitTimeOperator(c *TimeOperatorContext)

	// ExitNowOperator is called when exiting the nowOperator production.
	ExitNowOperator(c *NowOperatorContext)

	// ExitTime is called when exiting the time production.
	ExitTime(c *TimeContext)

	// ExitUnit is called when exiting the unit production.
	ExitUnit(c *UnitContext)

	// ExitInExpression is called when exiting the inExpression production.
	ExitInExpression(c *InExpressionContext)

	// ExitScope is called when exiting the scope production.
	ExitScope(c *ScopeContext)

	// ExitWithinExpression is called when exiting the withinExpression production.
	ExitWithinExpression(c *WithinExpressionContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)

	// ExitOperator is called when exiting the operator production.
	ExitOperator(c *OperatorContext)
}
