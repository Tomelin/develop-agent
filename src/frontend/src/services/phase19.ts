import { api } from "./api";
import { AdminQualityReport } from "@/types/phase19";

export const Phase19Service = {
  getAdminQualityReport: async (): Promise<AdminQualityReport> => {
    const response = await api.get("/admin/quality-report");
    return response.data;
  },
};

