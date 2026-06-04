# Projeto_Compilador_Faculdalde

**UNIC – Universidade de Cuiabá** **Relatório Técnico / Repositório Oficial** **Matéria:** Compiladores  
**Professor:** Edson Komatsu  
**Integrantes:** Emmanuel Duarte de Oliveira, Sandro Delmondes dos Anjos, Leandro Augusto Mestre Santana  
**Local/Ano:** Cuiabá/MT - 2026

---

## 1. Introdução
Este documento descreve a especificação formal e a implementação de um compilador didático desenvolvido para a disciplina de Compiladores. 

O sistema realiza o ciclo completo de tradução. Ele executa a **análise léxica** (identificação de tokens), **análise sintática** (construção da Árvore Sintática Abstrata - AST via *Top-Down*), **análise semântica** básica (verificação de declaração prévia de variáveis através de uma Tabela de Símbolos) e, por fim, a **geração de código alvo**, atuando como um transpilador que converte a AST validada em um código executável na linguagem **C**.

## 2. Arquitetura do Compilador e Diagrama de Fases
O compilador foi desenvolvido de forma modular em um único arquivo principal (`main.go`), estruturado nas seguintes fases e interfaces lógicas:

**Diagrama de Fases:**
`[Código Fonte (.txt)]` ➡️ **Léxico** ➡️ `[Tokens]` ➡️ **Sintático** ➡️ `[Árvore AST]` ➡️ **Semântico** ➡️ `[AST Validada]` ➡️ **Backend** ➡️ `[saida.c]` ➡️ **GCC** ➡️ `[Executável]`

* **Módulo Léxico:** Utiliza regras de expressões regulares mapeadas via `lexer.SimpleRule` para tokenização direta.
* **Módulo Sintático:** Mapeia a gramática EBNF diretamente para estruturas de dados (Structs) construindo a AST através da biblioteca Participle.
* **Módulo Semântico:** Funções recursivas validam o uso de variáveis cruzando dados com uma Tabela de Símbolos em tempo de execução. Bloqueia variáveis não declaradas.
* **Módulo de Backend (Geração de Código):** Módulo de transpilação que percorre a AST validada e traduz as instruções para um arquivo de saída formatado em C (`saida.c`), terceirizando a montagem do código de máquina final para o compilador GCC.

## 3. Justificativa das Ferramentas de Apoio
Em vez de utilizar geradores legados como Flex/Bison ou JavaCC, optou-se pela linguagem **Go** em conjunto com a biblioteca **Participle**. A escolha justifica-se pela capacidade do Participle de construir analisadores sintáticos descendentes (*Top-Down*) mapeando a gramática EBNF diretamente para as estruturas de dados nativas da linguagem. Isso garante uma arquitetura altamente modular e tipagem forte. Além disso, o binário gerado pelo Go é autossuficiente e de rápida execução, dispensando ambientes virtuais complexos como a JVM.

---

## 4. Gramática Formal (EBNF)
A sintaxe da linguagem segue o padrão ISO/IEC 14977. Abaixo, a representação das produções que compõem a Árvore Sintática Abstrata (AST):

```ebnf
(* Estrutura Principal *)
Programa = { Instrucao } ;
Instrucao = Atribuicao | Print | If | While ;

(* Regras de Produção *)
Atribuicao = Ident , "=" , Expressao ;
Print = "print" , "(" , Ident , ")" ;
If = "if" , "(" , Condicao , ")" , "{" , { Instrucao } , "}" , [ "else" , "{" , { Instrucao } , "}" ] ;
While = "while" , "(" , Condicao , ")" , "{" , { Instrucao } , "}" ;

Condicao = Termo , Operator , Termo ;
Expressao = Termo , [ Operator , Termo ] ;
Termo = Number | Ident ;

## 5. Especificação Léxica (Tokens)
* **Keyword:** `\b(if|else|while|print)\b` (Palavras reservadas para controle de fluxo e funções).
* **Ident:** `[a-zA-Z_][a-zA-Z0-9_]*` (Identificadores de variáveis).
* **Number:** `\d+` (Literais numéricos inteiros).
* **Operator:** `[=+*/><!-]` (Operadores de atribuição, aritméticos e relacionais).
* **Punct:** `[{}()]` (Delimitadores e pontuação para blocos).
* **Whitespace:** `\s+` (Espaços e quebras de linha - ignorados pelo sintático).

## 6. Regras de Sintaxe e Semântica
- **Estruturas de controle:** Suporta `if`/`else` e laços `while`.
- **Formatação:** Cada instrução deve estar em uma linha separada ou delimitada por chaves/blocos.
- **Tipos de dados:** O sistema lida primariamente com inteiros (`Number`).
- **Análise Semântica (Escopo):** Variáveis são instanciadas implicitamente na primeira atribuição. É estritamente proibido utilizar uma variável em operações matemáticas, lógicas ou em comandos `print` antes de sua atribuição. O compilador emitirá erro semântico nesses casos.

---

## 7. Configuração do Ambiente (Para Máquinas Novas)
Se você ou o avaliador estiverem executando este projeto em um computador novo, é necessário configurar os compiladores base:

### Passo 1: Instalar as Linguagens (Go e GCC)

**🖥️ Para Linux (Pop!_OS / Ubuntu):**
No terminal, execute os comandos:
```bash
sudo apt update
sudo apt install golang build-essential
🪟 Para Windows:

Instalar Go: Baixe e instale a versão mais recente diretamente do site oficial (go.dev).

Instalar GCC: Como o Windows não possui GCC nativo, você tem duas opções:

Opção A: Instalar o MinGW-w64 e adicioná-lo às variáveis de ambiente (PATH) do Windows.

Opção B (Recomendada): Rodar o projeto através do WSL (Windows Subsystem for Linux), o que permite usar exatamente os mesmos comandos de terminal do Linux.

Passo 2: Sincronizar as Bibliotecas (Crucial)
Abra o terminal (ou CMD/PowerShell) na pasta raiz do projeto e execute o comando abaixo. Isso fará o download da biblioteca participle utilizada na Árvore Sintática:

Bash
go mod tidy
Passo 3: Configurar o VS Code (Opcional)
Para uma melhor experiência de desenvolvimento e testes, recomendamos utilizar o Visual Studio Code com as extensões:

Go (Extensão oficial)

Code Runner (Para execução rápida com um clique)

8. Como Executar (Ambiente de Testes)
O compilador lê o código-fonte exclusivamente do arquivo teste.txt localizado na raiz do projeto.

Método 1: Script Automático (Linux/macOS/WSL)
Abra o arquivo teste.txt.

Apague o conteúdo explicativo e cole um dos exemplos da pasta /testes (códigos funcionando ou erros).

Salve o arquivo.

No terminal, execute:

Bash
./run.sh
Método 2: Extensão Code Runner (VS Code - Windows ou Linux)
Cole o código desejado no teste.txt e salve.

Abra o arquivo main.go.

Clique no botão de "Play" (Run Code) localizado no canto superior direito da tela.

Método 3: Execução Manual (Terminal, CMD ou PowerShell)
Se não quiser usar scripts, execute os passos diretamente:

Bash
# 1. Roda o compilador Go (Gera o arquivo saida.c)
go run main.go

# 2. Compila o arquivo C para código de máquina
gcc saida.c -o meu_programa

# 3. Executa o programa final
./meu_programa            # No Linux/Mac/WSL
meu_programa.exe          # No Windows
9. Estrutura da Pasta de Testes
Para facilitar a avaliação (testes de regressão e validação funcional), disponibilizamos uma pasta /testes com os seguintes cenários prontos:

codigo_funcionando_exemplo.txt: O caso de uso válido e completo.

erro_lexico.txt: Demonstra o tratamento de caracteres inválidos.

erro_sintatico.txt: Demonstra o tratamento de fechamentos de bloco incompletos.

erro_semantico.txt: Demonstra o bloqueio de variáveis fantasmas.

10. Referências Bibliográficas
MEDEIROS, João Antonio. Compiladores para Humanos. Disponível em: https://johnidm.gitbooks.io/compiladores-para-humanos/. Acesso em: jun. 2026.