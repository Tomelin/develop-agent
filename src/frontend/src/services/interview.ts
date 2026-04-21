import { api } from "./api";
import { InterviewSession } from "@/types/interview";

const getAuthToken = () =>
  typeof window === "undefined" ? null : localStorage.getItem("@agency:token");

const buildAuthHeaders = (): HeadersInit => {
  const token = getAuthToken();
  return {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
};

export const InterviewService = {
  getSession: async (projectId: string): Promise<InterviewSession> => {
    const response = await api.get(`/projects/${projectId}/interview`);
    return response.data;
  },

  regenerateVision: async (projectId: string): Promise<{ vision_markdown: string; generated_at?: string }> => {
    const response = await api.post(`/projects/${projectId}/interview/regenerate-vision`);
    return response.data;
  },

  confirmInterview: async (projectId: string): Promise<void> => {
    await api.post(`/projects/${projectId}/interview/confirm`);
  },

  streamMessage: async (
    projectId: string,
    content: string,
    handlers: {
      onToken: (token: string) => void;
      onDone: (payload?: { iteration_count?: number; message?: string }) => void;
      onError: (message: string) => void;
    },
  ) => {
    const baseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
    const response = await fetch(`${baseUrl}/projects/${projectId}/interview/message`, {
      method: "POST",
      headers: buildAuthHeaders(),
      credentials: "include",
      body: JSON.stringify({ content }),
    });

    if (!response.ok) {
      const rawError = await response.text();
      handlers.onError(rawError || "Falha ao enviar mensagem.");
      return;
    }

    if (!response.body) {
      handlers.onDone();
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = "";
    let streamDone = false;

    while (!streamDone) {
      const { done, value } = await reader.read();
      streamDone = done;
      buffer += decoder.decode(value || new Uint8Array(), { stream: !done });

      const chunks = buffer.split("\n\n");
      buffer = chunks.pop() || "";

      for (const chunk of chunks) {
        const lines = chunk.split("\n");
        let eventType = "message";
        const dataLines: string[] = [];

        for (const line of lines) {
          if (line.startsWith("event:")) {
            eventType = line.replace("event:", "").trim();
          }
          if (line.startsWith("data:")) {
            dataLines.push(line.replace("data:", "").trim());
          }
        }

        const raw = dataLines.join("\n");
        if (!raw) continue;

        if (eventType === "token" || eventType === "message") {
          try {
            const parsed = JSON.parse(raw) as { token?: string; content?: string };
            handlers.onToken(parsed.token || parsed.content || "");
          } catch {
            handlers.onToken(raw);
          }
        }

        if (eventType === "done") {
          try {
            handlers.onDone(JSON.parse(raw));
          } catch {
            handlers.onDone();
          }
        }

        if (eventType === "error") {
          try {
            const parsed = JSON.parse(raw) as { message?: string };
            handlers.onError(parsed.message || "Erro durante o streaming da resposta.");
          } catch {
            handlers.onError(raw);
          }
        }
      }
    }
  },
};
