import { Phase5CodeFile } from "./phase5";

export interface Phase13RunRequest {
  backend_base_url?: string;
  frontend_url?: string;
  include_devops: boolean;
}

export interface Phase13DeliveryReport {
  generated_at: string;
  project_id: string;
  project_status: string;
  include_devops: boolean;
  artifacts: string[];
  warnings?: string[];
}

export interface DeliveryArtifactGroups {
  documentation: Phase5CodeFile[];
  infrastructure: Phase5CodeFile[];
  summary?: Phase5CodeFile;
}
