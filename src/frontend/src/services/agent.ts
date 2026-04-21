import { api } from "./api";
import { AgentListResponse, Agent } from "../types/agent";

export const AgentService = {
  getAgents: async (page = 1, size = 10): Promise<AgentListResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      size: size.toString(),
    });
    const response = await api.get("/agents", { params });
    const payload = response.data;

    if (Array.isArray(payload)) {
      return {
        items: payload,
        total: payload.length,
        page,
        size,
        pages: 1,
      };
    }

    return {
      items: payload?.items ?? [],
      total: payload?.total ?? payload?.items?.length ?? 0,
      page: payload?.page ?? page,
      size: payload?.size ?? size,
      pages: payload?.pages ?? 1,
    };
  },

  getAgentById: async (id: string): Promise<Agent> => {
    const response = await api.get(`/agents/${id}`);
    return response.data;
  }
};
