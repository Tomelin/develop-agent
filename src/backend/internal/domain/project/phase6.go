package project

import "time"

type FunctionCoverage struct {
	Name    string  `json:"name"`
	Source  string  `json:"source"`
	Percent float64 `json:"percent"`
}

type PackageCoverage struct {
	Package string  `json:"package"`
	Percent float64 `json:"percent"`
}

type Phase6CoverageReport struct {
	GeneratedAt      time.Time          `json:"generated_at"`
	ThresholdPercent float64            `json:"threshold_percent"`
	TotalPercent     float64            `json:"total_percent"`
	Packages         []PackageCoverage  `json:"packages"`
	Functions        []FunctionCoverage `json:"functions"`
	RawGoTestOutput  string             `json:"raw_go_test_output,omitempty"`
}

type TestFailureKind string

const (
	TestFailureNone               TestFailureKind = "NONE"
	TestFailureTestImplementation TestFailureKind = "TEST_IMPLEMENTATION"
	TestFailureProjectBug         TestFailureKind = "PROJECT_BUG"
)

type Phase6ValidationResult struct {
	GoTestPassed       bool            `json:"go_test_passed"`
	GoRacePassed       bool            `json:"go_race_passed"`
	FrontendTestPassed bool            `json:"frontend_test_passed"`
	FailureKind        TestFailureKind `json:"failure_kind"`
	Details            string          `json:"details,omitempty"`
}
