package datastruct

type TestResult struct {
	Host       string
	Path       string
	Success    bool
	StatusCode int
	Error      error
}

type TestResultsMap map[string][]*TestResult
