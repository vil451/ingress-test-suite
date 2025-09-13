package test_loader

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"ingress-test-suite/internal/consts"
	ds "ingress-test-suite/internal/datastruct"
	networkingv1 "k8s.io/api/networking/v1"
)

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) LoadTestsFromDir(dir string) ([]*ds.IngressTestsFile, error) {
	entries, errReadDir := os.ReadDir(dir)
	if errReadDir != nil {
		return nil, errors.Wrapf(errReadDir, consts.ErrReadDir.Error(), dir)
	}

	result := make([]*ds.IngressTestsFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, errReadFile := os.ReadFile(path)
		if errReadFile != nil {
			return nil, errors.Wrapf(errReadFile, consts.ErrReadFile.Error(), path)
		}

		var tests ds.IngressTestsFile
		if err := json.Unmarshal(data, &tests); err != nil {
			return nil, errors.Wrapf(err, consts.ErrUnmarshalFile.Error(), path)
		}

		for _, t := range tests.Tests {
			if _, err := validatePathType(t.PathType); err != nil {
				return nil, err
			}
		}

		result = append(result, &tests)
	}

	return result, nil
}

func validatePathType(s string) (networkingv1.PathType, error) {
	switch s {
	case "Exact":
		pt := networkingv1.PathTypeExact
		return pt, nil
	case "Prefix":
		pt := networkingv1.PathTypePrefix
		return pt, nil
	case "ImplementationSpecific":
		pt := networkingv1.PathTypeImplementationSpecific
		return pt, nil
	default:
		return "", errors.Errorf(consts.ErrInvalidPathType.Error(), s)
	}
}
