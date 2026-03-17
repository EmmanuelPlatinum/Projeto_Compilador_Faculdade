package main

import (
	"fmt"
	"os"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// --- LÉXICO ---
var meuLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Whitespace", Pattern: `\s+`},
	{Name: "Keyword", Pattern: `\b(if|else|while|print)\b`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Number", Pattern: `\d+`},
	{Name: "Operator", Pattern: `[=+*/-]`},
	{Name: "Punct", Pattern: `[{}()]`},
})

// --- ESTRUTURA DA ÁRVORE SINTÁTICA (AST) ---
type Programa struct {
	Instrucoes []*Instrucao `parser:"{ @@ }"`
}

type Instrucao struct {
	Atribuicao *Atribuicao `parser:"( @@"`
	Print      *Print      `parser:"| @@"`
	If         *If         `parser:"| @@"`
	While      *While      `parser:"| @@ )"` // ADICIONADO WHILE
}

type Atribuicao struct {
	Variavel string     `parser:"@Ident '='"`
	Expressao *Expressao `parser:"@@"` // AGORA ACEITA EXPRESSÕES
}

type Expressao struct {
	Esquerda int    `parser:"@Number"`
	Op       string `parser:"[ @Operator"`
	Direita  int    `parser:"@Number ]"`
}

type Print struct {
	_     string `parser:"'print' '('"`
	Valor string `parser:"@Ident"`
	_     string `parser:"')'"`
}

type If struct {
	_        string       `parser:"'if' '('"`
	Condicao *Condicao    `parser:"@@ ')' '{'"`
	Corpo    []*Instrucao `parser:"{ @@ } '}'"`
}

type While struct {
	_        string       `parser:"'while' '('"`
	Condicao *Condicao    `parser:"@@ ')' '{'"`
	Corpo    []*Instrucao `parser:"{ @@ } '}'"`
}

type Condicao struct {
	Esquerda string `parser:"@Ident"`
	Op       string `parser:"@Operator"`
	Direita  int    `parser:"@Number"`
}

func main() {
	parser, err := participle.Build[Programa](participle.Lexer(meuLexer))
	if err != nil {
		panic(err)
	}

	dados, err := os.ReadFile("teste.txt")
	if err != nil {
		fmt.Println("Erro: Crie o arquivo 'teste.txt' na pasta!")
		return
	}

	ast := &Programa{}
	err = parser.ParseString("", string(dados), ast)
	if err != nil {
		fmt.Printf("❌ Erro de Sintaxe: %v\n", err)
		return
	}

	fmt.Println("✅ Sucesso! O compilador processou o While e as Expressões.")
}