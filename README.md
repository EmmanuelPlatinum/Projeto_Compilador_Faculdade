# Projeto_Compilador_Faculdalde

UNIC – Universidade de CuiabáRelatório Técnico: 

Especificação da Gramática do Analisador Léxico

Matéria: Compiladores

Professor: Edson Komatsu

Integrantes: Emmanuel Duarte de Oliveira, Sandro Delmondes dos Anjos, Leandro Augusto Mestre Santana

Cuiabá/MT 2026

# Compilador Simples em Go

## EBNF da Linguagem

```
Programa      ::= { Instrucao }
Instrução     ::= Atribuicao | Print | If | While
Atribuição    ::= Ident "=" Expressao
Expressão     ::= Termo [ Operator Termo ]
Termo         ::= Number | Ident
Print         ::= "print" "(" Ident ")"
If            ::= "if" "(" Condicao ")" "{" { Instrucao } "}" [ "else" "{" { Instrucao } "}" ]
While         ::= "while" "(" Condicao ")" "{" { Instrucao } "}"
Condição      ::= Termo Operator Termo
Ident         ::= [a-zA-Z_][a-zA-Z0-9_]*
Number        ::= [0-9]+
Operadores     ::= "=" | "+" | "-" | "*" | "/" | ">" | "<" | "!"
```

- **Tipos de dados**: Apenas inteiros (`Number`).
- **Operações**: Aritméticas e relacionais básicas.
- **Declaração de variáveis**: Implícita na atribuição.
- **Estruturas de controle**: `if`/`else`, `while`.
- **Operadores**: `=`, `+`, `-`, `*`, `/`, `>`, `<`, `!`

## Exemplo de Código (teste.txt)

```
x = 10
y = 20
if (x > y) {
    print(x)
} else {
    print(y)
}
while (x < 100) {
    x = x + 1
}
```

## Regras de Sintaxe
- Cada instrução deve estar em uma linha separada ou delimitada por chaves/blocos.
- Espaços em branco são ignorados.
- Variáveis são criadas na primeira atribuição.
- O compilador faz análise léxica, sintática e exibe tokens e erros léxicos.

## Execução
1. Compile com `go run main.go`.
2. O analisador léxico exibirá todos os tokens válidos e reportará erros léxicos, se houver.
3. O parser exibirá o número de instruções raiz encontradas.

1. IntroduçãoEste documento descreve a especificação formal da linguagem de programação desenvolvida para a disciplina de Compiladores. A implementação utiliza um analisador sintático descendente (Top-Down) baseado na biblioteca Participle para a linguagem Go. O sistema é capaz de processar instruções de atribuição, saída de dados e estruturas de controlo de fluxo.

2. Especificação Léxica (Tokens)

A análise léxica é definida por um conjunto de regras de expressão regular que identificam os símbolos básicos da linguagem:

Keyword: \b(if|else|while|print)\b

Descrição: Palavras reservadas para controle de fluxo e funções do sistema.

Ident: [a-zA-Z_][a-zA-Z0-9_]*

Descrição: Identificadores de variáveis (deve começar com letra ou sublinhado).

Number: \d+

Descrição: Literais numéricos inteiros (sequências de dígitos de 0 a 9).

Operator: [=+*/><!-]

Descrição: Operadores de atribuição, aritméticos e relacionais.

Punct: [{}()]

Descrição: Delimitadores e pontuação para blocos e expressões.

Whitespace: \s+

Descrição: Espaços, tabs e quebras de linha (identificados pelo léxico e ignorados pelo sintático).

3. Gramática Formal (EBNF)A sintaxe da linguagem segue o padrão ISO/IEC 14977. Abaixo, a representação das produções que compõem a Árvore Sintática Abstrata (AST):EBNF(* Estrutura Principal *)
Programa = { Instrucao } ;
Instrucao = Atribuicao
          | Print
          | If
          | While ;

(* Regras de Produção *)
Atribuição = Ident , "=" , Expressao ;
Print = "print" , "(" , Ident , ")" ;
If = "if" , "(" , Condicao , ")" , "{" , { Instrucao } , "}" , [ "else" , "{" , { Instrucao } , "}" ] ;
While = "while" , "(" , Condicao , ")" , "{" , { Instrucao } , "}" ;

Condição = Termo , Operator , Termo ;
Expressão = Termo , [ Operator , Termo ] ;
Termo = Number | Ident ;

4. Descrição das EstruturasPrograma: Raiz do parser, consistindo numa lista de instruções.If / Else: Estrutura condicional que permite a execução de blocos alternativos de código.While: Laço de repetição baseado numa condição lógica.Expressão: Suporta operações matemáticas simples entre números ou variáveis (Termos).Condição: Compara dois Termos (números ou variáveis) através de um operador lógico ou relacional.Termo: Unidade básica para cálculos e comparações, podendo ser um literal numérico ou um identificador (variável).

5. Exemplo de Implementação (Go)Abaixo, um fragmento do código que demonstra a integração do Lexer com o Parser:Go// Trecho do arquivo main.go
parser, err := participle.Build[Programa](
    participle.Lexer(meuLexer),
    participle.Elide("Whitespace"),
) 
// ... processamento do arquivo teste.txt