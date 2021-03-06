package rules

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"k8s.io/helm/pkg/lint/support"
	"os"
	"path/filepath"
	"text/template"
)

// Templates lints a chart's templates.
func Templates(linter *support.Linter) {
	templatespath := filepath.Join(linter.ChartDir, "templates")

	templatesExist := linter.RunLinterRule(support.WarningSev, validateTemplatesExistence(linter, templatespath))

	// Templates directory is optional for now
	if !templatesExist {
		return
	}

	linter.RunLinterRule(support.ErrorSev, validateTemplatesDir(linter, templatespath))
	linter.RunLinterRule(support.ErrorSev, validateTemplatesParseable(linter, templatespath))
}

func validateTemplatesExistence(linter *support.Linter, templatesPath string) (lintError support.LintError) {
	if _, err := os.Stat(templatesPath); err != nil {
		lintError = fmt.Errorf("Templates directory not found")
	}
	return
}

func validateTemplatesDir(linter *support.Linter, templatesPath string) (lintError support.LintError) {
	fi, err := os.Stat(templatesPath)
	if err == nil && !fi.IsDir() {
		lintError = fmt.Errorf("'templates' is not a directory")
	}
	return
}

func validateTemplatesParseable(linter *support.Linter, templatesPath string) (lintError support.LintError) {
	tpl := template.New("tpl").Funcs(sprig.TxtFuncMap())

	lintError = filepath.Walk(templatesPath, func(name string, fi os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if fi.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(name)
		if err != nil {
			lintError = fmt.Errorf("cannot read %s: %s", name, err)
			return lintError
		}

		newtpl, err := tpl.Parse(string(data))
		if err != nil {
			lintError = fmt.Errorf("error processing %s: %s", name, err)
			return lintError
		}
		tpl = newtpl
		return nil
	})

	return
}
