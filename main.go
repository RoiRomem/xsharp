package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
   TOKENIZER / LEXER SECTION
   --------------------------
   This section defines how to break the input source code into tokens.
   Each token is a meaningful unit like numbers, strings, identifiers, operators, etc.
*/

// Token struct holds the type, value, and location of each token.
type Token struct {
	Type   string // The type of token, e.g., "ID", "NUMBER", etc.
	Value  string // The literal value of the token.
	Line   int    // Line number where the token was found.
	Column int    // Column position in the line.
}

// tokenSpecs defines regex patterns for each type of token.
// Each entry has a token type and a regex that matches that token.
var tokenSpecs = []struct {
	Type  string
	Regex string
}{
	{"NUMBER", `\d+(\.\d*)?`},        // Integer or floating-point numbers.
	{"STRING", `"([^"\\]|\\.)*"`},    // Double-quoted strings with escapes.
	{"ID", `[A-Za-z_][A-Za-z0-9_]*`}, // Identifiers: names for variables, functions, etc.
	{"OP", `[+\-*/=<>!]`},            // Operators like +, -, *, /, etc.
	{"LPAREN", `\(`},                 // Left parenthesis.
	{"RPAREN", `\)`},                 // Right parenthesis.
	{"LBRACE", `{`},                  // Left brace.
	{"RBRACE", `}`},                  // Right brace.
	{"LANGLE", `<`},                  // Less-than sign.
	{"RANGLE", `>`},                  // Greater-than sign.
	{"COLON", `:`},                   // Colon, used in class inheritance.
	{"SEMICOLON", `;`},               // Semicolon, ends statements.
	{"COMMA", `,`},                   // Comma, separates parameters, etc.
	{"NEWLINE", `\n`},                // Newline characters.
	{"SKIP", `[ \t]+`},               // Skip over spaces and tabs.
	{"MISMATCH", `.`},                // Any other character (error if encountered).
}

// tokenize function scans the input code and produces a slice of Tokens.
func tokenize(code string) ([]Token, error) {
	var tokens []Token
	// Create a combined regex pattern for all token types.
	var patterns []string
	for _, spec := range tokenSpecs {
		// The regex is named with the token type.
		patterns = append(patterns, fmt.Sprintf("(?P<%s>%s)", spec.Type, spec.Regex))
	}
	regex := regexp.MustCompile(strings.Join(patterns, "|"))

	line := 1      // Current line number.
	lineStart := 0 // Position of the start of the current line.

	// Find all regex matches in the code.
	matches := regex.FindAllStringSubmatchIndex(code, -1)
	for _, match := range matches {
		// match[0] and match[1] are the start and end positions of the full match.
		fullStart, fullEnd := match[0], match[1]
		value := code[fullStart:fullEnd]
		var tokType string
		// Loop over each token spec to see which one matched.
		for i, spec := range tokenSpecs {
			// Each group index for token i is at positions 2*(i+1) and 2*(i+1)+1.
			start, end := match[2*(i+1)], match[2*(i+1)+1]
			if start != -1 && end != -1 {
				tokType = spec.Type
				break
			}
		}
		col := fullStart - lineStart // Calculate the column based on line start.
		switch tokType {
		case "SKIP":
			// Do nothing for spaces and tabs.
		case "NEWLINE":
			line++              // Increment line count.
			lineStart = fullEnd // Update the start position for the new line.
		case "MISMATCH":
			// Report an error for unrecognized characters.
			return nil, fmt.Errorf("unexpected token %q at line %d, col %d", value, line, col)
		default:
			// Append the token to our tokens slice.
			tokens = append(tokens, Token{Type: tokType, Value: value, Line: line, Column: col})
		}
	}
	// Append an "EOF" (end-of-file) token to signal the end of input.
	tokens = append(tokens, Token{Type: "EOF", Value: "", Line: line, Column: 0})
	return tokens, nil
}

/*
   ABSTRACT SYNTAX TREE (AST) SECTION
   ----------------------------------
   The AST represents the structure of your source code.
   We define different node types like Program, FunctionDecl, ClassDecl, etc.
*/

// Node interface: all AST nodes implement this.
type Node interface{}

// Program is the root node holding all top-level declarations.
type Program struct {
	Declarations []Node
}

// FunctionDecl represents a function declaration.
type FunctionDecl struct {
	RetType string  // Return type of the function.
	Name    string  // Function name.
	Params  []Param // Parameters of the function.
	Body    []Node  // Function body as a list of statements.
}

// Param represents a function parameter.
type Param struct {
	Type string // Parameter type.
	Name string // Parameter name.
}

// ClassDecl represents a class declaration.
type ClassDecl struct {
	Name    string // Class name.
	Parent  string // Parent class name, if any.
	Members []Node // Members: variables and functions.
}

// VarDecl represents a variable declaration.
type VarDecl struct {
	VarType string     // Variable type.
	Name    string     // Variable name.
	Default Expression // Default value (if provided).
}

// Expression represents a literal expression (number, string, or identifier).
type Expression struct {
	Value string // The literal value.
}

// Statement wraps an expression to be used as a statement.
type Statement struct {
	Expr Expression // The expression statement.
}

/*
   PARSER SECTION
   --------------
   The parser converts a stream of tokens into an AST.
   We implement a simple recursive descent parser to handle our language's grammar.
*/

type Parser struct {
	tokens []Token // All tokens from the lexer.
	pos    int     // Current position in the token slice.
}

// NewParser returns a new Parser instance.
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

// current returns the current token.
func (p *Parser) current() Token {
	return p.tokens[p.pos]
}

// consume moves to the next token and optionally checks the expected token type(s).
func (p *Parser) consume(expectedType ...string) Token {
	tok := p.current()
	if len(expectedType) > 0 {
		match := false
		for _, typ := range expectedType {
			// Allow matching against token type or literal value.
			if tok.Type == typ || tok.Value == typ {
				match = true
				break
			}
		}
		if !match {
			panic(fmt.Sprintf("Expected %v but got %s (%s) at line %d", expectedType, tok.Type, tok.Value, tok.Line))
		}
	}
	p.pos++
	return tok
}

// parse starts the parsing process and returns the Program AST node.
func (p *Parser) parse() Program {
	var decls []Node
	// Process tokens until we hit the EOF token.
	for p.current().Type != "EOF" {
		// If the token value is "class", parse a class declaration.
		if p.current().Value == "class" {
			decls = append(decls, p.parseClass())
		} else {
			decls = append(decls, p.parseFunction())
		}
	}
	return Program{Declarations: decls}
}

// parseFunction handles function declarations in the form:
// retType name ( params ) { body }
func (p *Parser) parseFunction() FunctionDecl {
	retType := p.consume("ID").Value // Function return type.
	name := p.consume("ID").Value    // Function name.
	p.consume("LPAREN")              // Consume '('.
	params := p.parseParams()        // Parse parameters.
	p.consume("RPAREN")              // Consume ')'.
	body := p.parseBlock()           // Parse function body enclosed in braces.
	return FunctionDecl{RetType: retType, Name: name, Params: params, Body: body}
}

// parseParams processes function parameters separated by commas.
func (p *Parser) parseParams() []Param {
	var params []Param
	// If the next token is RPAREN, there are no parameters.
	if p.current().Type == "RPAREN" {
		return params
	}
	// Loop until parameters are exhausted.
	for {
		paramType := p.consume("ID").Value // Parameter type.
		paramName := p.consume("ID").Value // Parameter name.
		params = append(params, Param{Type: paramType, Name: paramName})
		if p.current().Type == "COMMA" {
			p.consume("COMMA") // Consume comma between parameters.
		} else {
			break
		}
	}
	return params
}

// parseBlock processes a block of code enclosed in { }.
func (p *Parser) parseBlock() []Node {
	p.consume("LBRACE") // Consume '{'.
	var stmts []Node
	// Continue until the closing '}' is reached.
	for p.current().Type != "RBRACE" {
		stmts = append(stmts, p.parseStatement())
	}
	p.consume("RBRACE") // Consume '}'.
	return stmts
}

// parseStatement distinguishes between variable declarations and expression statements.
func (p *Parser) parseStatement() Node {
	// Lookahead: if we see two IDs in a row, assume it's a variable declaration.
	if p.current().Type == "ID" && p.tokens[p.pos+1].Type == "ID" {
		varType := p.consume("ID").Value // Variable type.
		varName := p.consume("ID").Value // Variable name.
		var def Expression               // Default value, if any.
		if p.current().Value == "=" {    // Check for assignment.
			p.consume("OP")           // Consume '=' operator.
			def = p.parseExpression() // Parse the default expression.
		}
		p.consume("SEMICOLON") // End of variable declaration.
		return VarDecl{VarType: varType, Name: varName, Default: def}
	}
	// Otherwise, parse an expression statement.
	expr := p.parseExpression()
	p.consume("SEMICOLON")
	return Statement{Expr: expr}
}

// parseExpression processes a simple literal expression.
func (p *Parser) parseExpression() Expression {
	tok := p.consume()
	// Support literals: NUMBER, STRING, or identifiers.
	if tok.Type == "NUMBER" || tok.Type == "STRING" || tok.Type == "ID" {
		return Expression{Value: tok.Value}
	}
	panic(fmt.Sprintf("Unexpected token in expression: %v", tok))
}

// parseClass handles class declarations in the form:
// class ClassName [: Parent] { members }
func (p *Parser) parseClass() ClassDecl {
	p.consume("ID")               // Consume the "class" keyword.
	name := p.consume("ID").Value // Class name.
	parent := ""
	// Optional inheritance: if a colon is present, read the parent class.
	if p.current().Type == "COLON" {
		p.consume("COLON")
		parent = p.consume("ID").Value
	}
	members := p.parseBlock() // Parse the class members enclosed in braces.
	return ClassDecl{Name: name, Parent: parent, Members: members}
}

/*
   CODE GENERATOR SECTION
   -----------------------
   The code generator traverses the AST and emits equivalent C code.
   It translates our custom language constructs into C constructs.
*/

type CodeGenerator struct {
	ast    Program         // The AST produced by the parser.
	code   strings.Builder // Used to build the output C code.
	indent string          // Current indentation string.
}

// NewCodeGenerator returns a new CodeGenerator.
func NewCodeGenerator(ast Program) *CodeGenerator {
	return &CodeGenerator{ast: ast, indent: ""}
}

// generate starts the code generation process.
func (cg *CodeGenerator) generate() string {
	cg.emitIncludes() // Emit standard C includes.
	// Process each top-level declaration.
	for _, decl := range cg.ast.Declarations {
		switch d := decl.(type) {
		case FunctionDecl:
			cg.emitFunction(d)
		case ClassDecl:
			cg.emitClass(d)
		}
	}
	return cg.code.String()
}

// emitIncludes writes the necessary C library includes.
func (cg *CodeGenerator) emitIncludes() {
	cg.code.WriteString("#include <stdio.h>\n#include <stdlib.h>\n#include <string.h>\n\n")
}

// emitFunction generates C code for a function declaration.
func (cg *CodeGenerator) emitFunction(fn FunctionDecl) {
	// Build parameter list as "type name" strings.
	var params []string
	for _, param := range fn.Params {
		params = append(params, fmt.Sprintf("%s %s", param.Type, param.Name))
	}
	// Emit function signature.
	cg.code.WriteString(fmt.Sprintf("%s %s(%s) {\n", fn.RetType, fn.Name, strings.Join(params, ", ")))
	cg.indent = "    " // Increase indentation for the function body.
	// Emit each statement in the function body.
	for _, stmt := range fn.Body {
		cg.emitStatement(stmt)
	}
	cg.code.WriteString("}\n\n") // Close the function.
}

// emitStatement generates C code for a single statement.
func (cg *CodeGenerator) emitStatement(stmt Node) {
	switch s := stmt.(type) {
	case VarDecl:
		// Variable declaration: type name [= default];
		line := fmt.Sprintf("%s%s %s", cg.indent, s.VarType, s.Name)
		if s.Default.Value != "" {
			// Check if the default is a number (no quotes needed) or string.
			if _, err := strconv.ParseFloat(s.Default.Value, 64); err == nil {
				line += " = " + s.Default.Value
			} else {
				line += " = " + s.Default.Value
			}
		}
		line += ";\n"
		cg.code.WriteString(line)
	case Statement:
		// Expression statement ends with a semicolon.
		cg.code.WriteString(fmt.Sprintf("%s%s;\n", cg.indent, s.Expr.Value))
	default:
		// Placeholder for any unhandled statements.
		cg.code.WriteString(fmt.Sprintf("%s// Unknown statement\n", cg.indent))
	}
}

// emitClass generates C code for a class declaration.
// It emits a C struct for the class and functions for its methods.
func (cg *CodeGenerator) emitClass(cls ClassDecl) {
	// Emit the struct definition for the class.
	cg.code.WriteString(fmt.Sprintf("typedef struct %s {\n", cls.Name))
	// For now, only handle member variable declarations.
	for _, mem := range cls.Members {
		if v, ok := mem.(VarDecl); ok {
			cg.code.WriteString(fmt.Sprintf("    %s %s;\n", v.VarType, v.Name))
		}
	}
	cg.code.WriteString(fmt.Sprintf("} %s;\n\n", cls.Name))
	// Emit methods as functions, with the first parameter being a pointer to the class instance.
	for _, mem := range cls.Members {
		if fn, ok := mem.(FunctionDecl); ok {
			params := []string{fmt.Sprintf("%s* this", cls.Name)}
			for _, param := range fn.Params {
				params = append(params, fmt.Sprintf("%s %s", param.Type, param.Name))
			}
			cg.code.WriteString(fmt.Sprintf("%s %s_%s(%s) {\n", fn.RetType, cls.Name, fn.Name, strings.Join(params, ", ")))
			cg.indent = "    "
			for _, stmt := range fn.Body {
				cg.emitStatement(stmt)
			}
			cg.code.WriteString("}\n\n")
		}
	}
}

/*
   MAIN FUNCTION
   -------------
   The entry point for the compiler. It ties together lexing, parsing, and code generation.
   It reads the input source file and writes the generated C code to the output file.
*/

func main() {
	// Ensure correct usage: compiler <input_file> <output_file>
	if len(os.Args) != 3 {
		fmt.Println("Usage: compiler <input_file> <output_file>")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	// Read the entire source code from the input file.
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}
	code := string(data)

	// --- Lexing ---
	tokens, err := tokenize(code)
	if err != nil {
		fmt.Println("Lexing error:", err)
		os.Exit(1)
	}

	// --- Parsing ---
	parser := NewParser(tokens)
	var ast Program
	// Catch any panic during parsing and report an error.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Parsing error:", r)
			os.Exit(1)
		}
	}()
	ast = parser.parse()

	// --- Code Generation ---
	gen := NewCodeGenerator(ast)
	cCode := gen.generate()

	// Write the generated C code to the output file.
	err = ioutil.WriteFile(outputFile, []byte(cCode), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		os.Exit(1)
	}
	fmt.Printf("C code generated and saved to %s\n", outputFile)
}
