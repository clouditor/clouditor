grammar CCL;

condition: assetType 'has' expression EOF;

assetType :
    simpleAssetType |
    filteredAssetType;

simpleAssetType: field;

filteredAssetType: field 'with' expression;

field : Identifier;
expression:
  simpleExpression |
  notExpression |
  inExpression;

simpleExpression:
  emptyExpression |
  withinExpression |
  comparison |
  '(' expression ')' ;

notExpression: 'not' expression;

emptyExpression: 'empty' field;

comparison: binaryComparison | timeComparison;

binaryComparison: field operator value;

timeComparison: field timeOperator (time unit | nowOperator);
timeOperator:
  BeforeOperator |
  AfterOperator |
  YoungerOperator |
  OlderOperator;
nowOperator: 'now';
time: Number;
unit:
  'seconds' |
  'days' |
  'months';

inExpression: simpleExpression 'in' scope field;

scope:
  'any' |
  'all';

withinExpression: field 'within' (value ','?)+;

value:
  StringLiteral |
  BooleanLiteral |
  Number;

operator:
  EqualsOperator |
  NotEqualsOperator |
  LessOrEqualsThanOperator |
  LessThanOperator |
  MoreThanOperator |
  MoreOrEqualsThanOperator |
  ContainsOperator;

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

BooleanLiteral:
  True |
  False;

True: 'true';
False: 'false';

Identifier
    : [a-zA-Z][a-zA-Z0-9.]*
    ;

Number: [0-9]+;

StringLiteral
    : '"' ~('"')* '"'
    ;

Whitespace
    : [ \t\u000C\r\n]+ -> skip
    ;
