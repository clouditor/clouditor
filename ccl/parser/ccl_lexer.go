// Code generated from CCL.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 35, 231,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23,
	9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9,
	28, 4, 29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33,
	4, 34, 9, 34, 3, 2, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3,
	5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3,
	7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3,
	9, 3, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 11, 3, 11,
	3, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 14, 3,
	14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 16, 3, 16, 3, 16,
	3, 17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 20, 3, 20, 3,
	21, 3, 21, 3, 21, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22,
	3, 22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 24, 3, 24, 3,
	24, 3, 24, 3, 24, 3, 24, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25,
	3, 25, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 27, 3, 27, 5, 27, 186,
	10, 27, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 29, 3, 29, 3, 29, 3, 29,
	3, 29, 3, 29, 3, 30, 3, 30, 7, 30, 201, 10, 30, 12, 30, 14, 30, 204, 11,
	30, 3, 31, 6, 31, 207, 10, 31, 13, 31, 14, 31, 208, 3, 32, 6, 32, 212,
	10, 32, 13, 32, 14, 32, 213, 3, 33, 3, 33, 7, 33, 218, 10, 33, 12, 33,
	14, 33, 221, 11, 33, 3, 33, 3, 33, 3, 34, 6, 34, 226, 10, 34, 13, 34, 14,
	34, 227, 3, 34, 3, 34, 2, 2, 35, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8,
	15, 9, 17, 10, 19, 11, 21, 12, 23, 13, 25, 14, 27, 15, 29, 16, 31, 17,
	33, 18, 35, 19, 37, 20, 39, 21, 41, 22, 43, 23, 45, 24, 47, 25, 49, 26,
	51, 27, 53, 28, 55, 29, 57, 30, 59, 31, 61, 32, 63, 33, 65, 34, 67, 35,
	3, 2, 8, 4, 2, 67, 92, 99, 124, 6, 2, 48, 48, 50, 59, 67, 92, 99, 124,
	3, 2, 50, 59, 4, 2, 48, 48, 50, 59, 3, 2, 36, 36, 5, 2, 11, 12, 14, 15,
	34, 34, 2, 236, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2,
	9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2,
	2, 17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2,
	2, 2, 25, 3, 2, 2, 2, 2, 27, 3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2,
	2, 2, 2, 33, 3, 2, 2, 2, 2, 35, 3, 2, 2, 2, 2, 37, 3, 2, 2, 2, 2, 39, 3,
	2, 2, 2, 2, 41, 3, 2, 2, 2, 2, 43, 3, 2, 2, 2, 2, 45, 3, 2, 2, 2, 2, 47,
	3, 2, 2, 2, 2, 49, 3, 2, 2, 2, 2, 51, 3, 2, 2, 2, 2, 53, 3, 2, 2, 2, 2,
	55, 3, 2, 2, 2, 2, 57, 3, 2, 2, 2, 2, 59, 3, 2, 2, 2, 2, 61, 3, 2, 2, 2,
	2, 63, 3, 2, 2, 2, 2, 65, 3, 2, 2, 2, 2, 67, 3, 2, 2, 2, 3, 69, 3, 2, 2,
	2, 5, 73, 3, 2, 2, 2, 7, 75, 3, 2, 2, 2, 9, 77, 3, 2, 2, 2, 11, 81, 3,
	2, 2, 2, 13, 87, 3, 2, 2, 2, 15, 91, 3, 2, 2, 2, 17, 99, 3, 2, 2, 2, 19,
	104, 3, 2, 2, 2, 21, 111, 3, 2, 2, 2, 23, 114, 3, 2, 2, 2, 25, 118, 3,
	2, 2, 2, 27, 122, 3, 2, 2, 2, 29, 129, 3, 2, 2, 2, 31, 131, 3, 2, 2, 2,
	33, 134, 3, 2, 2, 2, 35, 137, 3, 2, 2, 2, 37, 140, 3, 2, 2, 2, 39, 142,
	3, 2, 2, 2, 41, 144, 3, 2, 2, 2, 43, 147, 3, 2, 2, 2, 45, 156, 3, 2, 2,
	2, 47, 163, 3, 2, 2, 2, 49, 169, 3, 2, 2, 2, 51, 177, 3, 2, 2, 2, 53, 185,
	3, 2, 2, 2, 55, 187, 3, 2, 2, 2, 57, 192, 3, 2, 2, 2, 59, 198, 3, 2, 2,
	2, 61, 206, 3, 2, 2, 2, 63, 211, 3, 2, 2, 2, 65, 215, 3, 2, 2, 2, 67, 225,
	3, 2, 2, 2, 69, 70, 7, 106, 2, 2, 70, 71, 7, 99, 2, 2, 71, 72, 7, 117,
	2, 2, 72, 4, 3, 2, 2, 2, 73, 74, 7, 42, 2, 2, 74, 6, 3, 2, 2, 2, 75, 76,
	7, 43, 2, 2, 76, 8, 3, 2, 2, 2, 77, 78, 7, 112, 2, 2, 78, 79, 7, 113, 2,
	2, 79, 80, 7, 118, 2, 2, 80, 10, 3, 2, 2, 2, 81, 82, 7, 103, 2, 2, 82,
	83, 7, 111, 2, 2, 83, 84, 7, 114, 2, 2, 84, 85, 7, 118, 2, 2, 85, 86, 7,
	123, 2, 2, 86, 12, 3, 2, 2, 2, 87, 88, 7, 112, 2, 2, 88, 89, 7, 113, 2,
	2, 89, 90, 7, 121, 2, 2, 90, 14, 3, 2, 2, 2, 91, 92, 7, 117, 2, 2, 92,
	93, 7, 103, 2, 2, 93, 94, 7, 101, 2, 2, 94, 95, 7, 113, 2, 2, 95, 96, 7,
	112, 2, 2, 96, 97, 7, 102, 2, 2, 97, 98, 7, 117, 2, 2, 98, 16, 3, 2, 2,
	2, 99, 100, 7, 102, 2, 2, 100, 101, 7, 99, 2, 2, 101, 102, 7, 123, 2, 2,
	102, 103, 7, 117, 2, 2, 103, 18, 3, 2, 2, 2, 104, 105, 7, 111, 2, 2, 105,
	106, 7, 113, 2, 2, 106, 107, 7, 112, 2, 2, 107, 108, 7, 118, 2, 2, 108,
	109, 7, 106, 2, 2, 109, 110, 7, 117, 2, 2, 110, 20, 3, 2, 2, 2, 111, 112,
	7, 107, 2, 2, 112, 113, 7, 112, 2, 2, 113, 22, 3, 2, 2, 2, 114, 115, 7,
	99, 2, 2, 115, 116, 7, 112, 2, 2, 116, 117, 7, 123, 2, 2, 117, 24, 3, 2,
	2, 2, 118, 119, 7, 99, 2, 2, 119, 120, 7, 110, 2, 2, 120, 121, 7, 110,
	2, 2, 121, 26, 3, 2, 2, 2, 122, 123, 7, 121, 2, 2, 123, 124, 7, 107, 2,
	2, 124, 125, 7, 118, 2, 2, 125, 126, 7, 106, 2, 2, 126, 127, 7, 107, 2,
	2, 127, 128, 7, 112, 2, 2, 128, 28, 3, 2, 2, 2, 129, 130, 7, 46, 2, 2,
	130, 30, 3, 2, 2, 2, 131, 132, 7, 63, 2, 2, 132, 133, 7, 63, 2, 2, 133,
	32, 3, 2, 2, 2, 134, 135, 7, 35, 2, 2, 135, 136, 7, 63, 2, 2, 136, 34,
	3, 2, 2, 2, 137, 138, 7, 62, 2, 2, 138, 139, 7, 63, 2, 2, 139, 36, 3, 2,
	2, 2, 140, 141, 7, 62, 2, 2, 141, 38, 3, 2, 2, 2, 142, 143, 7, 64, 2, 2,
	143, 40, 3, 2, 2, 2, 144, 145, 7, 64, 2, 2, 145, 146, 7, 63, 2, 2, 146,
	42, 3, 2, 2, 2, 147, 148, 7, 101, 2, 2, 148, 149, 7, 113, 2, 2, 149, 150,
	7, 112, 2, 2, 150, 151, 7, 118, 2, 2, 151, 152, 7, 99, 2, 2, 152, 153,
	7, 107, 2, 2, 153, 154, 7, 112, 2, 2, 154, 155, 7, 117, 2, 2, 155, 44,
	3, 2, 2, 2, 156, 157, 7, 100, 2, 2, 157, 158, 7, 103, 2, 2, 158, 159, 7,
	104, 2, 2, 159, 160, 7, 113, 2, 2, 160, 161, 7, 116, 2, 2, 161, 162, 7,
	103, 2, 2, 162, 46, 3, 2, 2, 2, 163, 164, 7, 99, 2, 2, 164, 165, 7, 104,
	2, 2, 165, 166, 7, 118, 2, 2, 166, 167, 7, 103, 2, 2, 167, 168, 7, 116,
	2, 2, 168, 48, 3, 2, 2, 2, 169, 170, 7, 123, 2, 2, 170, 171, 7, 113, 2,
	2, 171, 172, 7, 119, 2, 2, 172, 173, 7, 112, 2, 2, 173, 174, 7, 105, 2,
	2, 174, 175, 7, 103, 2, 2, 175, 176, 7, 116, 2, 2, 176, 50, 3, 2, 2, 2,
	177, 178, 7, 113, 2, 2, 178, 179, 7, 110, 2, 2, 179, 180, 7, 102, 2, 2,
	180, 181, 7, 103, 2, 2, 181, 182, 7, 116, 2, 2, 182, 52, 3, 2, 2, 2, 183,
	186, 5, 55, 28, 2, 184, 186, 5, 57, 29, 2, 185, 183, 3, 2, 2, 2, 185, 184,
	3, 2, 2, 2, 186, 54, 3, 2, 2, 2, 187, 188, 7, 118, 2, 2, 188, 189, 7, 116,
	2, 2, 189, 190, 7, 119, 2, 2, 190, 191, 7, 103, 2, 2, 191, 56, 3, 2, 2,
	2, 192, 193, 7, 104, 2, 2, 193, 194, 7, 99, 2, 2, 194, 195, 7, 110, 2,
	2, 195, 196, 7, 117, 2, 2, 196, 197, 7, 103, 2, 2, 197, 58, 3, 2, 2, 2,
	198, 202, 9, 2, 2, 2, 199, 201, 9, 3, 2, 2, 200, 199, 3, 2, 2, 2, 201,
	204, 3, 2, 2, 2, 202, 200, 3, 2, 2, 2, 202, 203, 3, 2, 2, 2, 203, 60, 3,
	2, 2, 2, 204, 202, 3, 2, 2, 2, 205, 207, 9, 4, 2, 2, 206, 205, 3, 2, 2,
	2, 207, 208, 3, 2, 2, 2, 208, 206, 3, 2, 2, 2, 208, 209, 3, 2, 2, 2, 209,
	62, 3, 2, 2, 2, 210, 212, 9, 5, 2, 2, 211, 210, 3, 2, 2, 2, 212, 213, 3,
	2, 2, 2, 213, 211, 3, 2, 2, 2, 213, 214, 3, 2, 2, 2, 214, 64, 3, 2, 2,
	2, 215, 219, 7, 36, 2, 2, 216, 218, 10, 6, 2, 2, 217, 216, 3, 2, 2, 2,
	218, 221, 3, 2, 2, 2, 219, 217, 3, 2, 2, 2, 219, 220, 3, 2, 2, 2, 220,
	222, 3, 2, 2, 2, 221, 219, 3, 2, 2, 2, 222, 223, 7, 36, 2, 2, 223, 66,
	3, 2, 2, 2, 224, 226, 9, 7, 2, 2, 225, 224, 3, 2, 2, 2, 226, 227, 3, 2,
	2, 2, 227, 225, 3, 2, 2, 2, 227, 228, 3, 2, 2, 2, 228, 229, 3, 2, 2, 2,
	229, 230, 8, 34, 2, 2, 230, 68, 3, 2, 2, 2, 9, 2, 185, 202, 208, 213, 219,
	227, 3, 8, 2, 2,
}

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'has'", "'('", "')'", "'not'", "'empty'", "'now'", "'seconds'", "'days'",
	"'months'", "'in'", "'any'", "'all'", "'within'", "','", "'=='", "'!='",
	"'<='", "'<'", "'>'", "'>='", "'contains'", "'before'", "'after'", "'younger'",
	"'older'", "", "'true'", "'false'",
}

var lexerSymbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "EqualsOperator",
	"NotEqualsOperator", "LessOrEqualsThanOperator", "LessThanOperator", "MoreThanOperator",
	"MoreOrEqualsThanOperator", "ContainsOperator", "BeforeOperator", "AfterOperator",
	"YoungerOperator", "OlderOperator", "BooleanLiteral", "True", "False",
	"Identifier", "IntNumber", "FloatNumber", "StringLiteral", "Whitespace",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
	"T__9", "T__10", "T__11", "T__12", "T__13", "EqualsOperator", "NotEqualsOperator",
	"LessOrEqualsThanOperator", "LessThanOperator", "MoreThanOperator", "MoreOrEqualsThanOperator",
	"ContainsOperator", "BeforeOperator", "AfterOperator", "YoungerOperator",
	"OlderOperator", "BooleanLiteral", "True", "False", "Identifier", "IntNumber",
	"FloatNumber", "StringLiteral", "Whitespace",
}

type CCLLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

// NewCCLLexer produces a new lexer instance for the optional input antlr.CharStream.
//
// The *CCLLexer instance produced may be reused by calling the SetInputStream method.
// The initial lexer configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewCCLLexer(input antlr.CharStream) *CCLLexer {
	l := new(CCLLexer)
	lexerDeserializer := antlr.NewATNDeserializer(nil)
	lexerAtn := lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)
	lexerDecisionToDFA := make([]*antlr.DFA, len(lexerAtn.DecisionToState))
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "CCL.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// CCLLexer tokens.
const (
	CCLLexerT__0                     = 1
	CCLLexerT__1                     = 2
	CCLLexerT__2                     = 3
	CCLLexerT__3                     = 4
	CCLLexerT__4                     = 5
	CCLLexerT__5                     = 6
	CCLLexerT__6                     = 7
	CCLLexerT__7                     = 8
	CCLLexerT__8                     = 9
	CCLLexerT__9                     = 10
	CCLLexerT__10                    = 11
	CCLLexerT__11                    = 12
	CCLLexerT__12                    = 13
	CCLLexerT__13                    = 14
	CCLLexerEqualsOperator           = 15
	CCLLexerNotEqualsOperator        = 16
	CCLLexerLessOrEqualsThanOperator = 17
	CCLLexerLessThanOperator         = 18
	CCLLexerMoreThanOperator         = 19
	CCLLexerMoreOrEqualsThanOperator = 20
	CCLLexerContainsOperator         = 21
	CCLLexerBeforeOperator           = 22
	CCLLexerAfterOperator            = 23
	CCLLexerYoungerOperator          = 24
	CCLLexerOlderOperator            = 25
	CCLLexerBooleanLiteral           = 26
	CCLLexerTrue                     = 27
	CCLLexerFalse                    = 28
	CCLLexerIdentifier               = 29
	CCLLexerIntNumber                = 30
	CCLLexerFloatNumber              = 31
	CCLLexerStringLiteral            = 32
	CCLLexerWhitespace               = 33
)
