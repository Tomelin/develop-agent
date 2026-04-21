import { api } from "./api";
import {
  PromptGroup,
  PromptImportPayload,
  PromptListResponse,
  PromptPreviewResponse,
  PromptReorderItem,
  PromptTemplate,
  UserPrompt,
} from "@/types/prompt";

export const promptService = {
  async getPrompts(filters?: { group?: PromptGroup; enabled?: boolean }): Promise<PromptListResponse> {
    const response = await api.get("/prompts", { params: filters });
    return response.data;
  },

  async getPromptsByGroup(group: PromptGroup): Promise<UserPrompt[]> {
    const response = await api.get(`/prompts/${group}`);
    return response.data.items || response.data;
  },

  async createPrompt(payload: Omit<UserPrompt, "id" | "user_id" | "created_at" | "updated_at">): Promise<UserPrompt> {
    const response = await api.post("/prompts", payload);
    return response.data;
  },

  async updatePrompt(id: string, payload: Partial<UserPrompt>): Promise<UserPrompt> {
    const response = await api.put(`/prompts/${id}`, payload);
    return response.data;
  },

  async deletePrompt(id: string): Promise<void> {
    await api.delete(`/prompts/${id}`);
  },

  async reorderPrompts(items: PromptReorderItem[]): Promise<void> {
    await api.put("/prompts/reorder", { items });
  },

  async getPreview(group: PromptGroup, agentId?: string): Promise<PromptPreviewResponse> {
    const response = await api.get(`/prompts/preview/${group}`, {
      params: agentId ? { agent_id: agentId } : undefined,
    });
    return response.data;
  },

  async getTemplates(): Promise<PromptTemplate[]> {
    const response = await api.get("/prompts/templates");
    return response.data.items || response.data;
  },

  async createFromTemplate(templateId: string, group?: PromptGroup): Promise<UserPrompt> {
    const response = await api.post("/prompts/from-template", { template_id: templateId, group });
    return response.data;
  },

  async exportPrompts(): Promise<Blob> {
    const response = await api.get("/prompts/export", { responseType: "blob" });
    return response.data;
  },

  async importPrompts(payload: PromptImportPayload): Promise<{ imported: number; updated: number }> {
    const response = await api.post("/prompts/import", payload);
    return response.data;
  },
};
