package models

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type Finding struct {
	RuleName    string   `json:"rule_name"`
	Severity    Severity `json:"severity"`
	Description string   `json:"description"`
	FilePath    string   `json:"file_path"`
	LineNumber  int      `json:"line_number"`
	CodeSnippet string   `json:"code_snippet"`
}

type ReviewReport struct {
	TotalFindings int       `json:"total_findings"`
	Findings      []Finding `json:"findings"`
}
