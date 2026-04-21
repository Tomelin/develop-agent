export interface PackageCoverage {
  package: string;
  percent: number;
}

export interface FunctionCoverage {
  name: string;
  source: string;
  percent: number;
}

export interface Phase6CoverageReport {
  generated_at: string;
  threshold_percent: number;
  total_percent: number;
  packages: PackageCoverage[];
  functions: FunctionCoverage[];
  raw_go_test_output?: string;
}

export interface Phase6AnalyzeCoverageResponse {
  below_threshold: boolean;
  report: Phase6CoverageReport;
}

export type TestFailureKind = "NONE" | "TEST_IMPLEMENTATION" | "PROJECT_BUG";

export interface Phase6ValidationResult {
  go_test_passed: boolean;
  go_race_passed: boolean;
  frontend_test_passed: boolean;
  failure_kind: TestFailureKind;
  details?: string;
}
