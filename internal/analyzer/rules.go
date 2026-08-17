package analyzer

import (
	"regexp"

	"github.com/viniMSimoes/secpt-bot/internal/models"
)

type Rule struct {
	Name        string
	Severity    models.Severity
	Pattern     *regexp.Regexp
	Description string
}

var SecurityRules = []Rule{
	{
		Name:        "AWS Access Key ID",
		Severity:    models.SeverityCritical,
		Pattern:     regexp.MustCompile(`(?i)AKIA[0-9A-Z]{16}`),
		Description: "Chave de acesso da AWS detectada no código. Risco de comprometimento de infraestrutura em nuvem.",
	},
	{
		Name:        "Generic Secret / API Key",
		Severity:    models.SeverityHigh,
		Pattern:     regexp.MustCompile(`(?i)(api[_-]?key|secret[_-]?key|auth[_-]?token|private[_-]?key)\s*[:=]\s*["'][a-zA-Z0-9_\-]{8,}["']`),
		Description: "Possivel token de autenticação ou chave de API privada exposta em texto plano",
	},
	{
		Name:        "Hardcoded Password",
		Severity:    models.SeverityHigh,
		Pattern:     regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*["'][^"'\s]{6,}["']`),
		Description: "Senha fixa (hardcoded) detectada. Utilize variáveis de ambiente ou gerenciador de segredos.",
	},
	{
		Name:        "SQL Injection Risk (String Concatenation)",
		Severity:    models.SeverityHigh,
		Pattern:     regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE)\s+.*(\+|%s|\$|\.Format)`),
		Description: "Possível concatenação direta em query SQL. Risco de SQL Injection; utilize Prepared Statements / Queries parametrizadas.",
	},
	{
		Name:        "Dangerous Eval / Command Execution",
		Severity:    models.SeverityHigh,
		Pattern:     regexp.MustCompile(`(?i)(eval\(|exec\(|os\.system\(|exec\.Command\(.*request)`),
		Description: "Execução dinâmica de comandos ou código. Pode levar a Remote Code Execution (RCE).",
	},
}
