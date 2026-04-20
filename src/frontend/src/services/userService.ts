import { api } from './api';

export interface UserProfile {
  id: string;
  name: string;
  email: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export const userService = {
  getMe: async (): Promise<UserProfile> => {
    const response = await api.get<UserProfile>('/users/me');
    return response.data;
  },

  updateProfile: async (name: string): Promise<void> => {
    await api.put('/users/me', { name });
  },

  updatePassword: async (current_password: string, new_password: string): Promise<void> => {
    await api.put('/users/me/password', { current_password, new_password });
  }
};
