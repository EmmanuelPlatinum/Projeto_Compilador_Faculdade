#!/bin/bash
echo "=== INICIANDO COMPILADOR ==="
go run main.go

# Se o Go rodou sem erros (código de saída 0), ele compila o C e executa
if [ $? -eq 0 ]; then
    echo -e "\n=== SAÍDA DO PROGRAMA ==="
    gcc saida.c -o meu_programa
    ./meu_programa
else
    echo -e "\n⚠️ Compilação abortada devido a erros no código fonte."
fi

#Depois, no terminal, dê permissão de execução com o comando:
# chmod +x run.sh

# Agora,toda vez que você quiser testar o seu compilador, você só precisa digitar ./run.sh no terminal.
#Ele vai ler o código, gerar a árvore, transpilar, compilar no GCC e imprimir o resultado final em uma fração de segundos,
#tudo de uma vez

 


