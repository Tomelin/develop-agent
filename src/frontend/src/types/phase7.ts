export type SecuritySeverity = "CRITICAL" | "HIGH" | "MEDIUM" | "LOW";
export type SecurityFindingStatus = "OPEN" | "FIXED" | "ACCEPTED";

export interface SecurityFinding {
  id: string;
  title: string;
  category: string;
  severity: SecuritySeverity;
  cvss: number;
  cve?: string;
  description: string;
  poc?: string;
  file?: string;
  line?: number;
  remediation?: string;
  detected_by: string;
  status: SecurityFindingStatus;
}

export interface SecurityAuditSummary {
  score: number;
  critical_count: number;
  high_count: number;
  medium_count: number;
  low_count: number;
  total_findings: number;
}

export interface SecurityAutoRejectionResult {
  triggered: boolean;
  reason?: string;
  retry_count: number;
  returned_phase5: boolean;
  findings?: string[];
}

export interface SecurityAuditReport {
  generated_at: string;
  summary: SecurityAuditSummary;
  status: string;
  findings: SecurityFinding[];
  dependencies?: SecurityFinding[];
  static_analysis?: SecurityFinding[];
  auto_rejection: SecurityAutoRejectionResult;
}

export interface RunSecurityAuditPayload {
  backend_dir: string;
  frontend_dir: string;
  project_root_dir: string;
  high_retry_count: number;
}
