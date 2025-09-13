package internal

import (
	"context"
	"time"

	"github.com/pkg/errors"
	client_http "ingress-test-suite/internal/client/http"
	client_k8s "ingress-test-suite/internal/client/k8s"
	"ingress-test-suite/internal/consts"
	ds "ingress-test-suite/internal/datastruct"
	"ingress-test-suite/internal/ingress_manager"
	"ingress-test-suite/internal/test_loader"
)

type Logger interface {
	client_k8s.Logger
	client_http.Logger

	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type TestLoader interface {
	LoadTestsFromDir(dir string) ([]*ds.IngressTestsFile, error)
}

type ClientK8s interface {
	ingress_manager.K8sClient
}

type ClientHTTP interface {
	Test(entry *ds.IngressTestEntry) *ds.TestResult
}

type IngressManager interface {
	CheckExist(ctx context.Context, entry *ds.IngressTestEntry) (bool, error)
	Delete(ctx context.Context, entry *ds.IngressTestEntry) error
	Create(ctx context.Context, entry *ds.IngressTestEntry, c *ds.IngressTestsFile) error
}

type Tester struct {
	logger         Logger
	testLoader     TestLoader
	clientK8s      ClientK8s
	clientHttp     ClientHTTP
	ingressManager IngressManager
}

func NewTester(logger Logger) (*Tester, error) {
	tester := &Tester{
		logger:     logger,
		testLoader: test_loader.NewLoader(),
	}

	var err error
	tester.clientK8s, err = client_k8s.NewK8SClient(logger)
	if err != nil {
		return nil, errors.Wrap(err, consts.ErrMakeK8sClientSet.Error())
	}

	tester.ingressManager = ingress_manager.NewManager(tester.clientK8s)
	tester.clientHttp = client_http.NewHttpTester(tester.logger)

	return tester, nil
}

func (t *Tester) Run(ctx context.Context, testCasesPath string) (ds.TestResultsMap, error) {
	testCases, errLoadTests := t.testLoader.LoadTestsFromDir(testCasesPath)
	if errLoadTests != nil {
		return nil, errors.Wrap(errLoadTests, consts.ErrTesterRun.Error())
	}

	results := make(ds.TestResultsMap, len(testCases))

	for _, c := range testCases {
		t.logger.Printf(consts.MessageRunningIngressClass, c.IngressClassName)

		classResults := make([]*ds.TestResult, 0, len(c.Tests))
		for _, testLoad := range c.Tests {
			testResult := t.runSingleTest(ctx, c, &testLoad)
			classResults = append(classResults, testResult...)
		}

		if _, ok := results[c.IngressClassName]; ok {
			results[c.IngressClassName] = append(results[c.IngressClassName], classResults...)
		} else {
			results[c.IngressClassName] = classResults
		}

	}
	return results, nil
}

func (t *Tester) runSingleTest(ctx context.Context, file *ds.IngressTestsFile, entry *ds.IngressTestEntry) []*ds.TestResult {
	results := make([]*ds.TestResult, 0, 1)
	if entry.Create {
		if err := t.ensureIngressCreated(ctx, file, entry); err != nil {
			results = append(results, &ds.TestResult{
				Host:    entry.Host,
				Path:    entry.Path,
				Success: false,
				Error:   err,
			})
			return results
		}
		time.Sleep(2 * time.Second)
	}

	testRes := t.clientHttp.Test(entry)
	results = append(results, testRes)

	if entry.Create {
		if errCleanup := t.cleanupIngress(ctx, entry); errCleanup != nil {
			t.logger.Printf(errors.Wrapf(errCleanup, consts.ErrFailedIngressRuleDelete.Error(), entry.Host).Error())
		}
	}

	return results
}

func (t *Tester) ensureIngressCreated(ctx context.Context, file *ds.IngressTestsFile, entry *ds.IngressTestEntry) error {
	exists, err := t.ingressManager.CheckExist(ctx, entry)
	if err != nil {
		t.logger.Errorf(errors.Wrap(err, consts.MessageIngressRuleExist).Error())
		return errors.Wrap(err, consts.ErrFailedIngressRuleCreate.Error())
	}
	if exists {
		return nil
	}

	if errCreate := t.ingressManager.Create(ctx, entry, file); errCreate != nil {
		return errors.Wrap(errCreate, consts.ErrFailedIngressRuleCreate.Error())
	}
	return nil
}

func (t *Tester) cleanupIngress(ctx context.Context, entry *ds.IngressTestEntry) error {
	exists, err := t.ingressManager.CheckExist(ctx, entry)
	if err != nil {
		return t.ingressManager.Delete(ctx, entry)
	}
	if !exists {
		return errors.Errorf(consts.ErrIngressRuleNotExist.Error(), entry.Host)
	}

	return t.ingressManager.Delete(ctx, entry)
}
