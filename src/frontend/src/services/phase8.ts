import { isAxiosError } from "axios";
import { api } from "./api";
import {
  NotificationItem,
  PhaseArtifact,
  PhaseTrack,
  PhaseTrackStatus,
  TrackFeedbackItem,
  TriadTrackRuntime,
} from "@/types/phase8";


const toArray = <T>(payload: unknown): T[] => {
  if (Array.isArray(payload)) {
    return payload as T[];
  }

  if (payload && typeof payload === "object") {
    const record = payload as Record<string, unknown>;

    if (Array.isArray(record.items)) {
      return record.items as T[];
    }

    if (Array.isArray(record.rows)) {
      return record.rows as T[];
    }

    const objectValues = Object.values(record).filter((value) => value && typeof value === "object");
    if (objectValues.length) {
      return objectValues as T[];
    }
  }

  return [];
};

export const Phase8Service = {
  getTrackStatus: async (projectId: string, phaseNumber: number): Promise<PhaseTrackStatus[]> => {
    const response = await api.get(`/projects/${projectId}/phases/${phaseNumber}/tracks`);
    return toArray<PhaseTrackStatus>(response.data);
  },

  getArtifacts: async (projectId: string, phaseNumber: number): Promise<PhaseArtifact[]> => {
    try {
      const response = await api.get(`/projects/${projectId}/phases/${phaseNumber}/artifacts`);
      return toArray<PhaseArtifact>(response.data);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return [];
      }
      throw error;
    }
  },

  getTriadProgress: async (projectId: string, phaseNumber: number): Promise<TriadTrackRuntime[]> => {
    try {
      const response = await api.get(`/projects/${projectId}/phases/${phaseNumber}/triad-progress`);
      return toArray<TriadTrackRuntime>(response.data);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return [];
      }
      throw error;
    }
  },

  getFeedbackHistory: async (projectId: string, phaseNumber: number, track: PhaseTrack): Promise<TrackFeedbackItem[]> => {
    try {
      const response = await api.get(`/projects/${projectId}/phases/${phaseNumber}/tracks/${track}/feedbacks`);
      return toArray<TrackFeedbackItem>(response.data);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return [];
      }
      throw error;
    }
  },

  sendFeedback: async (projectId: string, phaseNumber: number, track: PhaseTrack, content: string): Promise<void> => {
    try {
      await api.post(`/projects/${projectId}/phases/${phaseNumber}/tracks/${track}/feedback`, { content });
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return;
      }
      throw error;
    }
  },

  approveTrack: async (projectId: string, phaseNumber: number, track: PhaseTrack): Promise<void> => {
    await api.post(`/projects/${projectId}/phases/${phaseNumber}/tracks/${track}/approve`);
  },

  getNotifications: async (): Promise<NotificationItem[]> => {
    try {
      const response = await api.get("/notifications");
      return toArray<NotificationItem>(response.data);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return [];
      }
      throw error;
    }
  },

  markNotificationAsRead: async (notificationId: string): Promise<void> => {
    try {
      await api.post(`/notifications/${notificationId}/read`);
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        return;
      }
      throw error;
    }
  },
};
