package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/viniMSimoes/secpt-bot/internal/analyzer"
	"github.com/viniMSimoes/secpt-bot/internal/github"
)

func main() {
	fmt.Println("🚀 Iniciando SecPR Security Scanner...")

	// 1. Obter o diff do Git (seja local ou no CI)
	baseBranch := os.Getenv("BASE_BRANCH")
	if baseBranch == "" {
		baseBranch = "HEAD~1" // Fallback para teste local
	}

	var cmd *exec.Cmd
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		// No CI, compara a branch atual com a base do PR
		cmd = exec.Command("git", "diff", "origin/"+baseBranch+"...HEAD")
	} else {
		// Localmente, compara com o commit anterior ou staged
		cmd = exec.Command("git", "diff", baseBranch)
	}

	diffOutput, err := cmd.Output()
	if err != nil {
		fmt.Printf("⚠️ Aviso ao obter git diff (%v). Tentando ler diff local...\n", err)
		cmd = exec.Command("git", "diff", "HEAD")
		diffOutput, _ = cmd.Output()
	}

	// 2. Escanear o diff
	scanner := analyzer.NewScanner()
	report := scanner.ScanDiff(string(diffOutput))
	markdownOutput := scanner.FormatMarkdown(report)

	// Imprimir no terminal local
	fmt.Printf("\n%s\n\n", markdownOutput)

	// 3. Se estiver rodando dentro do GitHub Action, posta o comentário no PR
	githubToken := os.Getenv("GITHUB_TOKEN")
	prNumber := os.Getenv("PR_NUMBER")
	repoOwner := os.Getenv("REPO_OWNER")
	repoName := os.Getenv("REPO_NAME")

	if githubToken != "" && prNumber != "" && repoOwner != "" && repoName != "" {
		fmt.Printf("📤 Enviando feedback para o PR #%s em %s/%s...\n", prNumber, repoOwner, repoName)
		client := github.NewClient(githubToken)
		err := client.PostPRComment(repoOwner, repoName, prNumber, markdownOutput)
		if err != nil {
			fmt.Printf("❌ Erro ao postar comentário no GitHub: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Comentário postado com sucesso no Pull Request!")
	} else {
		fmt.Println("ℹ️ Execução local finalizada (variáveis de CI não configuradas).")
	}

	// Se houver vulnerabilidades críticas, encerra com erro para falhar o build (opcional)
	if report.TotalFindings > 0 {
		fmt.Printf("⚠️ Foram encontrados %d problemas de segurança.\n", report.TotalFindings)
	}
}
