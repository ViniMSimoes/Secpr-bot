package analyzer

import (
	"fmt"
	"strings"

	"github.com/viniMSimoes/secpt-bot/internal/models"
)

type Scanner struct {
	rules []Rule
}

func NewScanner() *Scanner {
	return &Scanner{rules: SecurityRules}
}

func (s *Scanner) ScanDiff(diffContent string) models.ReviewReport {
	var findings []models.Finding
	lines := strings.Split(diffContent, "\n")

	currentFile := "unknown"
	currentLine := 0

	for _, line := range lines {
		if strings.HasPrefix(line, "+++ b/") {
			currentFile = strings.TrimPrefix(line, "+++ b/")
			currentLine = 0
			continue
		}

		if strings.HasPrefix(line, "@@") {
			continue
		}

		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			currentLine++
			codeContent := strings.TrimPrefix(line, "+")

			for _, rule := range s.rules {
				if rule.Pattern.MatchString(codeContent) {
					findings = append(findings, models.Finding{
						RuleName:    rule.Name,
						Severity:    rule.Severity,
						Description: rule.Description,
						FilePath:    currentFile,
						LineNumber:  currentLine,
						CodeSnippet: strings.TrimSpace(codeContent),
					})
				}
			}
		}
	}

	return models.ReviewReport{
		TotalFindings: len(findings),
		Findings:      findings,
	}
}

func (s *Scanner) FormatMarkdown(report models.ReviewReport) string {
	if report.TotalFindings == 0 {
		return "### 🛡️ **SecPR Bot — Relatório de Segurança**\n\n✅ **Nenhuma vulnerabilidade ou segredo exposto foi detectado nas alterações deste PR.** Ótimo trabalho!"
	}

	var sb strings.Builder
	sb.WriteString("### 🛡️ **SecPR Bot — Relatório de Segurança**\n\n")
	sb.WriteString(fmt.Sprintf("⚠️ **Atenção:** Foram encontrados **%d alerta(s) de segurança** nas alterações:\n\n", report.TotalFindings))
	sb.WriteString("| Severidade | Arquivo | Linha | Alerta | Detalhes |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")

	for _, f := range report.Findings {
		icon := "🟡"
		if f.Severity == models.SeverityCritical {
			icon = "🔴"
		} else if f.Severity == models.SeverityHigh {
			icon = "🟠"
		}

		sb.WriteString(fmt.Sprintf("| %s **%s** | `%s` | Linha %d | **%s** | %s |\n",
			icon, f.Severity, f.FilePath, f.LineNumber, f.RuleName, f.Description))
	}

	sb.WriteString("\n> 💡 *Dica de Remediação: Nunca versione credenciais no repositório. Mova chaves para o GitHub Secrets ou variáveis de ambiente e utilize queries parametrizadas.*")

	return sb.String()
}
