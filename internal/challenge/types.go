package challenge

type TestCase struct {
	Input       []interface{} `json:"input"`
	Expected    interface{}   `json:"expected"`
	Description string        `json:"description"`
}

type Challenge struct {
	Title          string     `json:"title"`
	Slug           string     `json:"slug"`
	Language       string     `json:"language"`
	Difficulty     string     `json:"difficulty"`
	FunctionName   string     `json:"functionName"`
	ParameterTypes []string   `json:"parameterTypes"`
	ReturnType     string     `json:"returnType"`
	Template       string     `json:"template"`
	TestCases      []TestCase `json:"testCases"`
	ConceptTags    []string   `json:"conceptTags"`
	TimeLimit      int        `json:"timeLimit"`
	MemoryLimit    int        `json:"memoryLimit"`
	Description    string     `json:"description,omitempty"`
}
