import { api } from "./api";
import { Phase14DeliveryReport, Phase14RunRequest } from "@/types/phase14";

export const Phase14Service = {
  run: async (projectId: string, payload: Phase14RunRequest): Promise<Phase14DeliveryReport> => {
    const response = await api.post(`/projects/${projectId}/phases/14/run`, payload);
    return response.data;
  },
};
