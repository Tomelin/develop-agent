import { api } from "./api";
import { RunSecurityAuditPayload, SecurityAuditReport } from "@/types/phase7";

export const Phase7Service = {
  runAudit: async (projectId: string, payload: RunSecurityAuditPayload): Promise<SecurityAuditReport> => {
    const response = await api.post(`/projects/${projectId}/phases/7/run-audit`, payload);
    return response.data;
  },
};
