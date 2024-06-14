package lexer

import "fmt"

type GojoToken struct {
	Type *GojoTokenType // The type of the token
	Text string         // The Text of the token
	Line int            // The line number of the token
}

func (t *GojoToken) String() string {
	return fmt.Sprintf("Token { Type: %s, Text: %q, Line: %d }", t.Type.StringInline(), t.Text, t.Line)
}

type GojoTokenType struct {
	Label      string // Name of the token type (e.g., "num")
	BeforeExpr bool   // Can be followed by an expression
	StartsExpr bool   // Can start an expression
	IsLoop     bool   // Is loop
}

func (t *GojoTokenType) String() string {
	return fmt.Sprintf("TokenType { Label: %s, BeforeExpr: %t, StartsExpr: %t, IsLoop: %t }",
		t.Label, t.BeforeExpr, t.StartsExpr, t.IsLoop)
}

// StringInline formats the GojoTokenType inline for better readability in nested structures.
func (t *GojoTokenType) StringInline() string {
	return fmt.Sprintf("Label: %s, BeforeExpr: %t, StartsExpr: %t, IsLoop: %t",
		t.Label, t.BeforeExpr, t.StartsExpr, t.IsLoop)
}

var TokenKeywords = map[string]*GojoTokenType{
	"break":      {Label: "break", BeforeExpr: true},
	"case":       {Label: "case", BeforeExpr: true},
	"catch":      {Label: "catch", BeforeExpr: true},
	"continue":   {Label: "continue", BeforeExpr: true},
	"debugger":   {Label: "debugger", BeforeExpr: true},
	"default":    {Label: "default", BeforeExpr: true},
	"do":         {Label: "do", IsLoop: true, BeforeExpr: true},
	"else":       {Label: "else", BeforeExpr: true},
	"finally":    {Label: "finally", BeforeExpr: true},
	"for":        {Label: "for", IsLoop: true},
	"function":   {Label: "function", StartsExpr: true},
	"if":         {Label: "if", BeforeExpr: true},
	"return":     {Label: "return", BeforeExpr: true},
	"switch":     {Label: "switch", BeforeExpr: true},
	"throw":      {Label: "throw", BeforeExpr: true},
	"try":        {Label: "try", BeforeExpr: true},
	"var":        {Label: "var", BeforeExpr: true},
	"let":        {Label: "var", BeforeExpr: true},
	"const":      {Label: "const", BeforeExpr: true},
	"while":      {Label: "while", IsLoop: true},
	"with":       {Label: "with", BeforeExpr: true},
	"new":        {Label: "new", BeforeExpr: true, StartsExpr: true},
	"this":       {Label: "this", StartsExpr: true},
	"super":      {Label: "super", StartsExpr: true},
	"class":      {Label: "class", StartsExpr: true},
	"extends":    {Label: "extends", BeforeExpr: true},
	"export":     {Label: "export", BeforeExpr: true},
	"import":     {Label: "import", StartsExpr: true},
	"null":       {Label: "null", StartsExpr: true},
	"true":       {Label: "true", StartsExpr: true},
	"false":      {Label: "false", StartsExpr: true},
	"in":         {Label: "in", BeforeExpr: true},
	"instanceof": {Label: "instanceof", BeforeExpr: true},
	"typeof":     {Label: "typeof", BeforeExpr: true, StartsExpr: true},
	"void":       {Label: "void", BeforeExpr: true, StartsExpr: true},
	"delete":     {Label: "delete", BeforeExpr: true, StartsExpr: true},
}

var TokenPunctuation = map[string]*GojoTokenType{
	"(":   {Label: "(", BeforeExpr: true, StartsExpr: true},
	")":   {Label: ")", BeforeExpr: false},
	"{":   {Label: "{", BeforeExpr: true, StartsExpr: true},
	"}":   {Label: "}", BeforeExpr: false},
	"[":   {Label: "[", BeforeExpr: true, StartsExpr: true},
	"]":   {Label: "]", BeforeExpr: false},
	",":   {Label: ",", BeforeExpr: false},
	";":   {Label: ";", BeforeExpr: false},
	":":   {Label: ":", BeforeExpr: false},
	".":   {Label: ".", BeforeExpr: true, StartsExpr: true},
	"?":   {Label: "?", BeforeExpr: true, StartsExpr: true},
	"=>":  {Label: "=>", BeforeExpr: true},
	"...": {Label: "...", BeforeExpr: true},
}

var TokenOperators = map[string]*GojoTokenType{
	"=":   {Label: "=", BeforeExpr: true},
	"+":   {Label: "+", BeforeExpr: true, StartsExpr: true}, // Can be unary or binary
	"-":   {Label: "-", BeforeExpr: true, StartsExpr: true}, // Can be unary or binary
	"*":   {Label: "*", BeforeExpr: true, StartsExpr: true},
	"/":   {Label: "/", BeforeExpr: true, StartsExpr: true},
	"!":   {Label: "!", BeforeExpr: true, StartsExpr: true},
	"~":   {Label: "~", BeforeExpr: true, StartsExpr: true},
	"==":  {Label: "==", BeforeExpr: true},  // Equality
	"!=":  {Label: "!=", BeforeExpr: true},  // Equality
	"!==": {Label: "!==", BeforeExpr: true}, // Equality
	"===": {Label: "===", BeforeExpr: true}, // Equality
	"<":   {Label: "<", BeforeExpr: true},   // Relational
	">":   {Label: ">", BeforeExpr: true},   // Relational
	"<=":  {Label: "<=", BeforeExpr: true},  // Relational
	">=":  {Label: ">=", BeforeExpr: true},  // Relational
	"&&":  {Label: "&&", BeforeExpr: true},  // Logical AND
	"||":  {Label: "||", BeforeExpr: true},  // Logical OR
	"|":   {Label: "|", BeforeExpr: true},   // Bitwise OR
	"^":   {Label: "^", BeforeExpr: true},   // Bitwise XOR
	"&":   {Label: "&", BeforeExpr: true},   // Bitwise AND
	"<<":  {Label: "<<", BeforeExpr: true},  // Bit Shift
	">>":  {Label: ">>", BeforeExpr: true},  // Bit Shift
	">>>": {Label: ">>>", BeforeExpr: true}, // Bit Shift
	"%":   {Label: "%", BeforeExpr: true},   // Modulo
	"**":  {Label: "**", BeforeExpr: true},
	"??":  {Label: "??", BeforeExpr: true}, // Coalesce
}

var TokenText = map[string]*GojoTokenType{
	"identifier": {Label: "identifier", StartsExpr: true},
	"eof":        {Label: "eof"},
}

var TokenLiterals = map[string]*GojoTokenType{
	"number":   {Label: "number", StartsExpr: true},
	"string":   {Label: "string", StartsExpr: true},
	"regexp":   {Label: "regexp", StartsExpr: true},
	"template": {Label: "template", StartsExpr: true},
}
