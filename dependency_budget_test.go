package xdatly

import (
	"os/exec"
	"strings"
	"testing"
)

func TestDependencyBudget_NoDatlyImplementationImports(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go list -deps failed: %v\n%s", err, string(output))
	}
	for _, line := range strings.Split(string(output), "\n") {
		importPath := strings.TrimSpace(line)
		if importPath == "" {
			continue
		}
		if strings.Contains(importPath, "github.com/viant/datly_1") || strings.Contains(importPath, "github.com/viant/datly/") || importPath == "github.com/viant/datly" {
			t.Fatalf("dependency budget violated by import %s", importPath)
		}
	}
}
