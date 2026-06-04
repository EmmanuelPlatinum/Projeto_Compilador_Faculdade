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

	// ... seu código que já estava lá em cima (parser, etc) ...

    fmt.Printf("✅ Sucesso! O compilador encontrou %d instrução(ões) raiz no código.\n", len(ast.Instrucoes))

    // === EXECUÇÃO DA FASE 5 (SEMÂNTICA) ===
    fmt.Println("\n🚀 Iniciando Análise Semântica...")
    
    // Cria a Tabela de Símbolos (um mapa para guardar o nome das variáveis)
    tabelaSimbolos := make(map[string]bool)
    
    // Roda a análise
    errSemantico := analisarSemantica(ast.Instrucoes, tabelaSimbolos)
    if errSemantico != nil {
        fmt.Printf("❌ Erro Semântico: %v\n", errSemantico)
        return // Trava a compilação aqui se der erro
    }
    
    fmt.Println("✅ Análise Semântica aprovada! Nenhuma variável usada sem declaração.")

    // === EXECUÇÃO DA FASE 6 (GERAÇÃO DE CÓDIGO) ===
    fmt.Println("\n⚙️ Gerando código executável (C)...")
    errGen := gerarCodigoC(ast.Instrucoes, tabelaSimbolos)
    if errGen != nil {
        fmt.Printf("❌ Erro na geração de código: %v\n", errGen)
        return
    }
    fmt.Println("✅ Código gerado com sucesso no arquivo 'saida.c'!")

} // 🚨 A FUNÇÃO MAIN() TERMINA EXATAMENTE AQUI! NENHUM COMANDO DEVE FICAR ABAIXO DESSA CHAVE!


// =============================================================================
// TODAS AS FUNÇÕES EXTRAS FICAM AQUI PARA BAIXO (FORA DO MAIN)
// =============================================================================

// --- FASE 5: ANÁLISE SEMÂNTICA ---
func analisarSemantica(instrucoes []*Instrucao, simbolos map[string]bool) error {
	for _, inst := range instrucoes {
		if inst.Atribuicao != nil {
			if err := checarExpressao(inst.Atribuicao.Expressao, simbolos); err != nil {
				return err
			}
			simbolos[inst.Atribuicao.Variavel] = true
		} else if inst.Print != nil {
			if !simbolos[inst.Print.Valor] {
				return fmt.Errorf("variável '%s' usada no print sem ser declarada", inst.Print.Valor)
			}
		} else if inst.If != nil {
			if err := checarCondicao(inst.If.Condicao, simbolos); err != nil {
				return err
			}
			if err := analisarSemantica(inst.If.Corpo, simbolos); err != nil {
				return err
			}
			if inst.If.Else != nil {
				if err := analisarSemantica(inst.If.Else, simbolos); err != nil {
					return err
				}
			}
		} else if inst.While != nil {
			if err := checarCondicao(inst.While.Condicao, simbolos); err != nil {
				return err
			}
			if err := analisarSemantica(inst.While.Corpo, simbolos); err != nil {
				return err
			}
		}
	}
	return nil
}

func checarExpressao(exp *Expressao, simbolos map[string]bool) error {
	if exp.Esquerda.Variavel != nil && !simbolos[*exp.Esquerda.Variavel] {
		return fmt.Errorf("variável '%s' não declarada", *exp.Esquerda.Variavel)
	}
	if exp.Direita != nil && exp.Direita.Variavel != nil && !simbolos[*exp.Direita.Variavel] {
		return fmt.Errorf("variável '%s' não declarada", *exp.Direita.Variavel)
	}
	return nil
}

func checarCondicao(cond *Condicao, simbolos map[string]bool) error {
	if cond.Esquerda.Variavel != nil && !simbolos[*cond.Esquerda.Variavel] {
		return fmt.Errorf("variável '%s' não declarada na condição", *cond.Esquerda.Variavel)
	}
	if cond.Direita.Variavel != nil && !simbolos[*cond.Direita.Variavel] {
		return fmt.Errorf("variável '%s' não declarada na condição", *cond.Direita.Variavel)
	}
	return nil
}

// --- FASE 6: GERAÇÃO DE CÓDIGO (C) ---
func gerarCodigoC(instrucoes []*Instrucao, simbolos map[string]bool) error {
	file, err := os.Create("saida.c")
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("#include <stdio.h>\n\n")
	file.WriteString("int main() {\n")

	if len(simbolos) > 0 {
		file.WriteString("    int ")
		var vars []string
		for v := range simbolos {
			vars = append(vars, v)
		}
		file.WriteString(strings.Join(vars, ", ") + ";\n")
	}

	gerarInstrucoesC(file, instrucoes, 1)

	file.WriteString("    return 0;\n")
	file.WriteString("}\n")
	return nil
}

func gerarInstrucoesC(file *os.File, instrucoes []*Instrucao, nivel int) {
	ident := strings.Repeat("    ", nivel)
	for _, inst := range instrucoes {
		if inst.Atribuicao != nil {
			esq := termoString(inst.Atribuicao.Expressao.Esquerda)
			dir := ""
			if inst.Atribuicao.Expressao.Direita != nil {
				dir = " " + inst.Atribuicao.Expressao.Op + " " + termoString(inst.Atribuicao.Expressao.Direita)
			}
			file.WriteString(fmt.Sprintf("%s%s = %s%s;\n", ident, inst.Atribuicao.Variavel, esq, dir))
		} else if inst.Print != nil {
			file.WriteString(fmt.Sprintf("%sprintf(\"%%d\\n\", %s);\n", ident, inst.Print.Valor))
		} else if inst.If != nil {
			condEsq := termoString(inst.If.Condicao.Esquerda)
			condOp := inst.If.Condicao.Op
			if condOp == "=" { condOp = "==" } 
			condDir := termoString(inst.If.Condicao.Direita)

			file.WriteString(fmt.Sprintf("%sif (%s %s %s) {\n", ident, condEsq, condOp, condDir))
			gerarInstrucoesC(file, inst.If.Corpo, nivel+1)
			if inst.If.Else != nil {
				file.WriteString(fmt.Sprintf("%s} else {\n", ident))
				gerarInstrucoesC(file, inst.If.Else, nivel+1)
			}
			file.WriteString(fmt.Sprintf("%s}\n", ident))
		} else if inst.While != nil {
			condEsq := termoString(inst.While.Condicao.Esquerda)
			condOp := inst.While.Condicao.Op
			if condOp == "=" { condOp = "==" } 
			condDir := termoString(inst.While.Condicao.Direita)

			file.WriteString(fmt.Sprintf("%swhile (%s %s %s) {\n", ident, condEsq, condOp, condDir))
			gerarInstrucoesC(file, inst.While.Corpo, nivel+1)
			file.WriteString(fmt.Sprintf("%s}\n", ident))
		}
	}
}

func termoString(t *Termo) string {
	if t.Numero != nil {
		return fmt.Sprintf("%d", *t.Numero)
	}
	return *t.Variavel
}