package scale

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
)

// FrameworkFileName is the expected framework metadata file in a catalog
// directory.
const FrameworkFileName = "framework.json"

// DomainsDirName is the expected subdirectory holding per-domain files.
const DomainsDirName = "domains"

// ExternalDirName is the expected subdirectory holding codified external
// maturity models.
const ExternalDirName = "external"

// ParseFramework parses a Framework document from JSON.
func ParseFramework(data []byte) (*Framework, error) {
	var f Framework
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parsing framework: %w", err)
	}
	return &f, nil
}

// ParseDomain parses a Domain document from JSON.
func ParseDomain(data []byte) (*Domain, error) {
	var d Domain
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("parsing domain: %w", err)
	}
	return &d, nil
}

// ParseAssessment parses an Assessment document from JSON.
func ParseAssessment(data []byte) (*Assessment, error) {
	var a Assessment
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("parsing assessment: %w", err)
	}
	return &a, nil
}

// LoadFrameworkFS assembles and validates a Framework from a catalog
// directory in fsys: root/framework.json for framework metadata and
// narratives, plus root/domains/*.json appended in filename order. Domains
// already inlined in framework.json are kept and extended.
func LoadFrameworkFS(fsys fs.FS, root string) (*Framework, error) {
	data, err := fs.ReadFile(fsys, path.Join(root, FrameworkFileName))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", FrameworkFileName, err)
	}
	f, err := ParseFramework(data)
	if err != nil {
		return nil, err
	}

	domainFiles, err := listJSONFiles(fsys, path.Join(root, DomainsDirName))
	if err != nil {
		return nil, err
	}
	for _, name := range domainFiles {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("reading domain file %s: %w", name, err)
		}
		d, err := ParseDomain(data)
		if err != nil {
			return nil, fmt.Errorf("domain file %s: %w", name, err)
		}
		f.Domains = append(f.Domains, *d)
	}

	externalFiles, err := listJSONFiles(fsys, path.Join(root, ExternalDirName))
	if err != nil {
		return nil, err
	}
	for _, name := range externalFiles {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("reading external model file %s: %w", name, err)
		}
		em, err := ParseExternalModel(data)
		if err != nil {
			return nil, fmt.Errorf("external model file %s: %w", name, err)
		}
		f.ExternalModels = append(f.ExternalModels, *em)
	}

	if err := f.Validate(); err != nil {
		return nil, fmt.Errorf("framework invalid: %w", err)
	}
	return f, nil
}

// ParseExternalModel parses an ExternalModel document from JSON.
func ParseExternalModel(data []byte) (*ExternalModel, error) {
	var em ExternalModel
	if err := json.Unmarshal(data, &em); err != nil {
		return nil, fmt.Errorf("parsing external model: %w", err)
	}
	return &em, nil
}

// listJSONFiles returns the .json files directly inside dir, sorted by name,
// with paths joined to dir. A missing directory yields no files and no error.
func listJSONFiles(fsys fs.FS, dir string) ([]string, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || path.Ext(e.Name()) != ".json" {
			continue
		}
		names = append(names, path.Join(dir, e.Name()))
	}
	sort.Strings(names)
	return names, nil
}

// LoadFrameworkDir assembles and validates a Framework from a catalog
// directory on disk.
func LoadFrameworkDir(dir string) (*Framework, error) {
	return LoadFrameworkFS(os.DirFS(dir), ".")
}

// LoadAssessmentFile reads and parses an Assessment from disk. Pass a
// non-nil framework to also validate observations against the catalog.
func LoadAssessmentFile(filename string, f *Framework) (*Assessment, error) {
	data, err := os.ReadFile(filename) //nolint:gosec // G304: path provided by caller
	if err != nil {
		return nil, fmt.Errorf("reading assessment %s: %w", filename, err)
	}
	a, err := ParseAssessment(data)
	if err != nil {
		return nil, err
	}
	if f != nil {
		if err := a.Validate(f); err != nil {
			return nil, fmt.Errorf("assessment invalid: %w", err)
		}
	}
	return a, nil
}
