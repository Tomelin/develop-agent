import { api } from "@/services/api";
import { PlatformStatusResponse } from "@/types/status";

export const statusService = {
  getPlatformStatus: async (): Promise<PlatformStatusResponse> => {
    const response = await api.get("/status");
    return response.data;
  },
};
