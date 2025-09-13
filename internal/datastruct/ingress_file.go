package datastruct

type IngressTestEntry struct {
	Host           string `json:"host"`
	Path           string `json:"path"`
	Service        string `json:"service"`
	PathType       string `json:"pathType"`
	ExpectedStatus int    `json:"expectedStatus"`
	Namespace      string `json:"namespace"`
	Port           int    `json:"port"`
	ExtPort        int    `json:"extPort"`
	Create         bool   `json:"create"`
}

type IngressTestsFile struct {
	IngressClassName string             `json:"ingressClassName"`
	Tests            []IngressTestEntry `json:"tests"`
}
