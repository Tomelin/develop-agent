import { api } from "./api";
import { Phase13DeliveryReport, Phase13RunRequest } from "@/types/phase13";

export const Phase13Service = {
  run: async (projectId: string, payload: Phase13RunRequest): Promise<Phase13DeliveryReport> => {
    const response = await api.post(`/projects/${projectId}/phases/13/run`, payload);
    return response.data;
  },
};
