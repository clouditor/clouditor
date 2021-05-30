grammar CCL;

condition: assetType 'has' expression EOF;

assetType: simpleAssetType | filteredAssetType;

simpleAssetType: field;

filteredAssetType: field 'with' expression;

field: Identifier;
expression: simpleExpression | notExpression | inExpression;

simpleExpression:
	isEmptyExpression
	| withinExpression
	| comparison
	| '(' expression ')';

notExpression: 'not' expression;

isEmptyExpression: 'empty' field;

comparison: binaryComparison | timeComparison;

binaryComparison: field operator value;

timeComparison: field timeOperator (time unit | nowOperator);
timeOperator:
	BeforeOperator
	| AfterOperator
	| YoungerOperator
	| OlderOperator;
nowOperator: 'now';
time: IntNumber;
unit: 'seconds' | 'days' | 'months';

inExpression: simpleExpression 'in' scope field;

scope: 'any' | 'all';

withinExpression: field 'within' (value ','?)+;

value: StringLiteral | BooleanLiteral | IntNumber | FloatNumber;

operator:
	EqualsOperator
	| NotEqualsOperator
	| LessOrEqualsThanOperator
	| LessThanOperator
	| MoreThanOperator
	| MoreOrEqualsThanOperator
	| ContainsOperator;

EqualsOperator: '==';
NotEqualsOperator: '!=';
LessOrEqualsThanOperator: '<=';
LessThanOperator: '<';
MoreThanOperator: '>';
MoreOrEqualsThanOperator: '>=';
ContainsOperator: 'contains';

BeforeOperator: 'before';
AfterOperator: 'after';
YoungerOperator: 'younger';
OlderOperator: 'older';

BooleanLiteral: True | False;

True: 'true';
False: 'false';

Identifier: [a-zA-Z][a-zA-Z0-9.]*;

IntNumber: [0-9]+;
FloatNumber: [0-9.]+;

StringLiteral: '"' ~('"')* '"';

Whitespace: [ \t\u000C\r\n]+ -> skip;
