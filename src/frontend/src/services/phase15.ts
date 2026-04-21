import { api } from "./api";
import { MarketingChannel, MarketingWebhookResult, Phase15DeliveryReport, Phase15RunRequest } from "@/types/phase15";

export const Phase15Service = {
  run: async (projectId: string, payload: Phase15RunRequest): Promise<Phase15DeliveryReport> => {
    const response = await api.post(`/projects/${projectId}/phases/15/run`, payload);
    return response.data;
  },

  downloadPack: async (projectId: string, channels: MarketingChannel[]): Promise<{ blob: Blob; pieces: number; filename?: string }> => {
    const params = new URLSearchParams();
    if (channels.length > 0) {
      params.set("channels", channels.join(","));
    }

    const response = await api.get(`/projects/${projectId}/marketing/export`, {
      params,
      responseType: "blob",
    });

    const pieces = Number(response.headers["x-marketing-pieces"] ?? 0);
    const disposition = String(response.headers["content-disposition"] ?? "");
    const filename = disposition.match(/filename="?([^\"]+)"?/)?.[1];

    return { blob: response.data, pieces, filename };
  },

  configureWebhook: async (projectId: string, url: string): Promise<MarketingWebhookResult> => {
    const response = await api.post(`/projects/${projectId}/marketing/webhooks`, { url });
    return response.data;
  },
};
