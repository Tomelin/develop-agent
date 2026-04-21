import { api } from "./api";
import {
  BillingGroupedItem,
  BillingPricingTable,
  BillingProjectDetails,
  BillingRecord,
  BillingSummary,
} from "@/types/billing";

export interface BillingQueryParams {
  from?: string;
  to?: string;
  page?: number;
  limit?: number;
}

const sanitizeParams = (params?: BillingQueryParams) => {
  if (!params) return undefined;

  return {
    from: params.from,
    to: params.to,
    page: params.page ?? 1,
    limit: params.limit ?? 50,
  };
};

export const BillingService = {
  async getPricing(): Promise<BillingPricingTable> {
    const response = await api.get("/billing/pricing");
    return response.data;
  },

  async getSummary(params?: BillingQueryParams): Promise<BillingSummary> {
    const response = await api.get("/billing/summary", {
      params: sanitizeParams(params),
    });
    return response.data;
  },

  async getByModel(params?: BillingQueryParams): Promise<BillingGroupedItem[]> {
    const response = await api.get("/billing/by-model", {
      params: sanitizeParams(params),
    });
    return response.data.items || [];
  },

  async getByPhase(params?: BillingQueryParams): Promise<BillingGroupedItem[]> {
    const response = await api.get("/billing/by-phase", {
      params: sanitizeParams(params),
    });
    return response.data.items || [];
  },

  async getTopProjects(params?: BillingQueryParams): Promise<BillingGroupedItem[]> {
    const response = await api.get("/billing/top-projects", {
      params: sanitizeParams(params),
    });
    return response.data.items || [];
  },

  async getProjectBilling(projectId: string, params?: BillingQueryParams): Promise<BillingProjectDetails> {
    const response = await api.get(`/projects/${projectId}/billing`, {
      params: sanitizeParams(params),
    });
    return response.data;
  },

  async getRecords(params?: BillingQueryParams & { project_id?: string; provider?: string }): Promise<BillingRecord[]> {
    const response = await api.get("/billing/export", {
      params: {
        ...sanitizeParams(params),
        project_id: params?.project_id,
        provider: params?.provider,
        format: "json",
      },
    });
    return response.data || [];
  },

  async exportBilling(
    format: "csv" | "json",
    params?: BillingQueryParams & { project_id?: string; provider?: string },
  ): Promise<Blob> {
    const response = await api.get("/billing/export", {
      params: {
        ...sanitizeParams(params),
        project_id: params?.project_id,
        provider: params?.provider,
        format,
      },
      responseType: "blob",
    });

    return response.data;
  },

  async updateProjectBudget(projectId: string, budgetUsd: number): Promise<void> {
    await api.put(`/projects/${projectId}/budget`, { budget_usd: budgetUsd });
  },
};
