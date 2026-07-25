// Command scale validates SCALE catalogs and renders story reports.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	scale "github.com/ProductBuildersHQ/scale"
	"github.com/ProductBuildersHQ/scale/report"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		usage()
		return fmt.Errorf("a subcommand is required")
	}
	switch args[0] {
	case "validate":
		return runValidate(args[1:])
	case "report":
		return runReport(args[1:])
	case "export":
		return runExport(args[1:])
	case "help", "-h", "--help":
		usage()
		return nil
	default:
		usage()
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `Usage:
  scale validate -catalog <dir> [-assessment <file>]
  scale report   -catalog <dir> -assessment <file> [-prev <file>] [-model <id>] -o <out.html>
  scale export   -catalog <dir> -o <framework.json>

report renders the SCALE aspect story report by default; -model renders the
same catalog and assessment through an external maturity model's ladder
(e.g. -model aws-observability-maturity).

export assembles framework.json + domains/*.json + external/*.json into a
single validated framework JSON IR, e.g. for the <scale-report> web component.
`)
}

func runValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	catalogDir := fs.String("catalog", "catalog", "catalog directory (framework.json, domains/, external/)")
	assessmentFile := fs.String("assessment", "", "optional assessment file to validate against the catalog")
	if err := fs.Parse(args); err != nil {
		return err
	}

	f, err := scale.LoadFrameworkDir(*catalogDir)
	if err != nil {
		return err
	}
	fmt.Printf("catalog OK: %d domains, %d external models\n", len(f.Domains), len(f.ExternalModels))

	if *assessmentFile != "" {
		a, err := scale.LoadAssessmentFile(*assessmentFile, f)
		if err != nil {
			return err
		}
		fmt.Printf("assessment OK: period %s, %d observations\n", a.Period, len(a.Observations))
	}
	return nil
}

func runExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	catalogDir := fs.String("catalog", "catalog", "catalog directory")
	outFile := fs.String("o", "framework.json", "output JSON file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	f, err := scale.LoadFrameworkDir(*catalogDir)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling framework: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(*outFile, data, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", *outFile, err)
	}
	fmt.Printf("wrote %s (%d domains, %d external models)\n", *outFile, len(f.Domains), len(f.ExternalModels))
	return nil
}

func runReport(args []string) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	catalogDir := fs.String("catalog", "catalog", "catalog directory")
	assessmentFile := fs.String("assessment", "", "current-period assessment file (required)")
	prevFile := fs.String("prev", "", "optional prior-period assessment file for deltas and movers")
	modelID := fs.String("model", "", "render through an external maturity model's ladder (model ID) instead of the aspect report")
	outFile := fs.String("o", "scale-report.html", "output HTML file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *assessmentFile == "" {
		return fmt.Errorf("-assessment is required")
	}

	f, err := scale.LoadFrameworkDir(*catalogDir)
	if err != nil {
		return err
	}
	curr, err := scale.LoadAssessmentFile(*assessmentFile, f)
	if err != nil {
		return err
	}
	opts := &report.Options{GeneratedAt: time.Now().Format("2006-01-02")}
	if *prevFile != "" {
		prev, err := scale.LoadAssessmentFile(*prevFile, f)
		if err != nil {
			return err
		}
		opts.Prev = prev
	}

	var html []byte
	if *modelID != "" {
		html, err = report.ModelHTML(f, curr, *modelID, opts)
	} else {
		html, err = report.HTML(f, curr, opts)
	}
	if err != nil {
		return err
	}
	if err := os.WriteFile(*outFile, html, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", *outFile, err)
	}
	fmt.Printf("wrote %s\n", *outFile)
	return nil
}
