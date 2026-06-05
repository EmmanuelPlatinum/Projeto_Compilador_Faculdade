package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Descobre a pasta atual para evitar erros de caminho
	dirAtual, _ := os.Getwd()
	
	// Configuração padrão (assumindo que rodou da pasta raiz do projeto)
	caminhoHTML := "web/index.html"
	caminhoTeste := "teste.txt"
	comandoBash := "./run.sh"

	// Se o usuário rodou o comando de dentro da pasta 'web', o servidor ajusta as rotas sozinho:
	if filepath.Base(dirAtual) == "web" {
		caminhoHTML = "index.html"
		caminhoTeste = "../teste.txt"
		comandoBash = "cd .. && ./run.sh"
	}

	// 1. Rota para mostrar a interface visual (Página Web)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, caminhoHTML)
	})

	// 2. Rota de API que recebe o código do site e manda pro compilador
	http.HandleFunc("/executar", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		codigo := r.FormValue("codigo")

		// Salva o código recebido no teste.txt correto
		os.WriteFile(caminhoTeste, []byte(codigo), 0644)

		// Executa o compilador
		cmd := exec.Command("bash", "-c", comandoBash)
		out, err := cmd.CombinedOutput()

		// Devolve a resposta do terminal pro Front-End
		if err != nil {
			fmt.Fprintf(w, "Erro ao compilar:\n%v\n\nLogs do Terminal:\n%s", err, string(out))
			return
		}
		fmt.Fprintf(w, "%s", string(out))
	})

	fmt.Println("🚀 Servidor Web rodando! Abra seu navegador em: http://localhost:8080")
	fmt.Println("⚠️  Para desligar o servidor depois, aperte CTRL + C no terminal.")
	
	// Inicia o servidor na porta 8080
	http.ListenAndServe(":8080", nil)
}