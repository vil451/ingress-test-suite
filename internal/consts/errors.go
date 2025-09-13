package consts

import "github.com/pkg/errors"

var (
	ErrTesterRun                   = errors.New("failed to run tester")
	ErrTestCasesDirPath            = errors.New("TESTS_PATH is not set (either as an env variable or CLI argument)")
	ErrReadDir                     = errors.New("failed to read dir %s")
	ErrReadFile                    = errors.New("failed to read file %s")
	ErrUnmarshalFile               = errors.New("failed to unmarshal file %s")
	ErrInvalidPathType             = errors.New("invalid pathType: %s")
	ErrMakeK8sClientSet            = errors.New("failed to create k8s client set")
	ErrLoadKubeConfig              = errors.New("could not determine kubeConfig")
	ErrBuildKubeConfig             = errors.New("failed to build config from kubeconfig")
	ErrFailedCreateKubeClient      = errors.New("failed to create k8s client")
	ErrFailedConvertPathType       = errors.New("failed to convert pathType")
	ErrFailedCheckExistIngressRule = errors.New("failed to check ingress %s")
	ErrHttpRequestFailed           = errors.New("http request failed")
	ErrFailedCloseResponseBody     = errors.New("warning: failed to close response body")
	ErrIngressRuleNotExist         = errors.New("ingress rule does not exist: %s")
	ErrFailedIngressRuleCreate     = errors.New("failed to create ingress")
	ErrFailedIngressRuleDelete     = errors.New("failed to delete ingress for host=%s")
)
