import { api } from "./api";
import {
  NotificationItem,
  PhaseArtifact,
  PhaseTrack,
  PhaseTrackStatus,
  TrackFeedbackItem,
  TriadTrackRuntime,
} from "@/types/phase8";

export const Phase8Service = {
  getTrackStatus: async (projectId: string, phaseNumber: number): Promise<PhaseTrackStatus[]> => {
    const response = await api.get(`/projects/${projectId}/phases/${phaseNumber}/tracks`);
    return response.data?.items ?? response.data ?? [];
  },

  getArtifacts: async (projectId: string, phaseNumber: number): Promise<PhaseArtifact[]> => {
    const response = await api.get(`/projects/${projectId}/phases/${phaseNumber}/artifacts`);
    return response.data?.items ?? response.data ?? [];
  },

  getTriadProgress: async (projectId: string, phaseNumber: number): Promise<TriadTrackRuntime[]> => {
    const response = await api.get(`/projects/${projectId}/phases/${phaseNumber}/triad-progress`);
    return response.data?.items ?? response.data ?? [];
  },

  getFeedbackHistory: async (projectId: string, phaseNumber: number, track: PhaseTrack): Promise<TrackFeedbackItem[]> => {
    const response = await api.get(`/projects/${projectId}/phases/${phaseNumber}/tracks/${track}/feedbacks`);
    return response.data?.items ?? response.data ?? [];
  },

  sendFeedback: async (projectId: string, phaseNumber: number, track: PhaseTrack, content: string): Promise<void> => {
    await api.post(`/projects/${projectId}/phases/${phaseNumber}/tracks/${track}/feedback`, { content });
  },

  approveTrack: async (projectId: string, phaseNumber: number, track: PhaseTrack): Promise<void> => {
    await api.post(`/projects/${projectId}/phases/${phaseNumber}/tracks/${track}/approve`);
  },

  getNotifications: async (): Promise<NotificationItem[]> => {
    const response = await api.get("/notifications");
    return response.data?.items ?? response.data ?? [];
  },

  markNotificationAsRead: async (notificationId: string): Promise<void> => {
    await api.post(`/notifications/${notificationId}/read`);
  },
};
