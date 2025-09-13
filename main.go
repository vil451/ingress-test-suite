package main

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"ingress-test-suite/internal"
	"ingress-test-suite/internal/consts"
	"ingress-test-suite/internal/pkg/logger"
)

var version = "dev"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		cancel()
	}()

	log := logger.New(logger.InfoLevel)

	testCasesPath := parseFlags(log)

	tester, err := internal.NewTester(log)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	results, runErr := tester.Run(ctx, testCasesPath)
	if runErr != nil {
		if errors.Is(runErr, consts.ErrMakeK8sClientSet) {
			log.Fatalf(runErr.Error())
			log.Exit(31)
			return
		}

		log.Fatalf(runErr.Error())
		log.Exit(1)
	}

	exitCode := 0
	for caseName, caseResults := range results {
		log.Infof(consts.MessageResultTestCaseInfo, caseName)
		for _, r := range caseResults {
			if r.Success {
				log.Infof(consts.MessageStatus, r.Host, r.Path, consts.MessageOK, r.StatusCode)
			} else {
				log.Errorf(consts.MessageStatus, r.Host, r.Path, consts.MessageFail, r.StatusCode)
				exitCode = 12
			}
		}
	}

	log.Exit(exitCode)
}

func parseFlags(log *logger.Logger) string {
	testCasesPath := ""
	testCasesPathEnv := os.Getenv("TESTS_PATH")
	pflag.StringVarP(&testCasesPath,
		"tests-path", "t", testCasesPathEnv, consts.MessageTestCasesPathDescription)
	versionFlag := pflag.BoolP("version", "v", false, consts.MessageVersionDescription)
	helpFlag := pflag.BoolP("help", "h", false, consts.MessageHelp)
	pflag.Parse()

	if *helpFlag {
		log.Infof(consts.MessageUsage)
		pflag.PrintDefaults()
		log.Exit(0)
	}

	if *versionFlag {
		log.Infof(consts.MessageVersion, version)
	}

	testCasesPath = strings.TrimSpace(testCasesPath)
	if testCasesPath == "" {
		log.Fatalf(consts.ErrTestCasesDirPath.Error())
	}
	return testCasesPath
}
