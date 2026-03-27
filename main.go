package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// --- FASE 1: LÉXICO ---
var meuLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Whitespace", Pattern: `\s+`},
	{Name: "Keyword", Pattern: `\b(if|else|while|print)\b`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Number", Pattern: `\d+`},
	{Name: "Operator", Pattern: `[=+*/><!-]`}, // Adicionado <, >, !
	{Name: "Punct", Pattern: `[{}()]`},
})

// --- FASE 2: SINTÁTICO (AST) ---
type Programa struct {
	Instrucoes []*Instrucao `parser:"{ @@ }"`
}

type Instrucao struct {
	Atribuicao *Atribuicao `parser:"( @@"`
	Print      *Print      `parser:"| @@"`
	If         *If         `parser:"| @@"`
	While      *While      `parser:"| @@ )"`
}

type Atribuicao struct {
	Variavel  string     `parser:"@Ident '='"`
	Expressao *Expressao `parser:"@@"`
}

type Expressao struct {
	Esquerda *Termo `parser:"@@"`
	Op       string `parser:"[ @Operator"`
	Direita  *Termo `parser:"@@ ]"`
}

type Termo struct {
	Numero   *int    `parser:"( @Number"`
	Variavel *string `parser:"| @Ident )"`
}

type Print struct {
	Valor string `parser:"'print' '(' @Ident ')'"`
}

type If struct {
	Condicao *Condicao    `parser:"'if' '(' @@ ')' '{'"`
	Corpo    []*Instrucao `parser:"{ @@ } '}'"`
	Else     []*Instrucao `parser:"( 'else' '{' { @@ } '}' )?"`
}

type While struct {
	Condicao *Condicao    `parser:"'while' '(' @@ ')' '{'"`
	Corpo    []*Instrucao `parser:"{ @@ } '}'"`
}

type Condicao struct {
	Esquerda *Termo `parser:"@@"`
	Op       string `parser:"@Operator"`
	Direita  *Termo `parser:"@@"`
}

// --- FASE 3: EXECUÇÃO E TESTE ---
func main() {
	// Lendo o arquivo teste.txt
	dados, err := os.ReadFile("teste.txt")
	if err != nil {
		fmt.Println("❌ Erro: Arquivo 'teste.txt' não encontrado na mesma pasta do main.go.")
		return
	}

	fmt.Printf("🔍 O compilador leu o seguinte texto do arquivo:\n[%s]\n\n", string(dados))

	// --- ANÁLISE LÉXICA MANUAL ---
	// Usa io.Reader corretamente
	lex, err := meuLexer.Lex("teste.txt", strings.NewReader(string(dados)))
	if err != nil {
		fmt.Printf("❌ Erro ao inicializar o Lexer: %v\n", err)
		return
	}

	// 1. Descobre o ID do token Whitespace para ignorá-lo depois
	whitespaceID := meuLexer.Symbols()["Whitespace"]

	// 2. CORREÇÃO: Criamos o nosso próprio mapa reverso manualmente!
	mapaReverso := make(map[lexer.TokenType]string)
	for nome, id := range meuLexer.Symbols() {
		mapaReverso[id] = nome
	}

	fmt.Println("🔍 Tokens encontrados:")
	for {
		token, err := lex.Next()
		if err != nil {
			// Erro léxico: mostra linha e coluna se possível
			fmt.Printf("❌ Erro Léxico: %v\n", err)
			return
		}
		if token.EOF() {
			break
		}
		// Ignora espaços em branco usando o ID correto
		if token.Type == whitespaceID {
			continue
		}
		// CORREÇÃO: Usamos o nosso mapaReverso em vez da função inventada pelo Copilot
		fmt.Printf("Tipo: %-10s Valor: %-10q (linha %d, coluna %d)\n",
			mapaReverso[token.Type], token.Value, token.Pos.Line, token.Pos.Column)
	}

	// --- ANÁLISE SINTÁTICA ---
	parser, err := participle.Build[Programa](
		participle.Lexer(meuLexer),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		panic(err)
	}

	ast, err := parser.ParseString("teste.txt", string(dados))
	if err != nil {
		// Erro sintático: mostra linha e coluna se possível
		if parseErr, ok := err.(interface{ Position() lexer.Position }); ok {
			pos := parseErr.Position()
			fmt.Printf("❌ Erro de Sintaxe na linha %d, coluna %d: %v\n", pos.Line, pos.Column, err)
		} else {
			fmt.Printf("❌ Erro de Sintaxe: %v\n", err)
		}
		return
	}

	fmt.Printf("✅ Sucesso! O compilador encontrou %d instrução(ões) raiz no código.\n", len(ast.Instrucoes))
}