import { httpClient } from '../../api/client';
import { ApiResponse } from '../../api/types';
import { Group, CreateAndUpdateGroupPayload } from './types';

export const groupApi = {
  getGroups: async (): Promise<ApiResponse<Group[]>> => {
    try {
      const response = await httpClient.get<{ data: Group[] }>('/groups');
      return { success: true, data: response.data.data };
    } catch (error: any) {
      return { success: false, error: error.message || 'Failed to fetch groups' };
    }
  },

  createGroup: async (payload: CreateAndUpdateGroupPayload): Promise<ApiResponse<Group>> => {
    try {
      const response = await httpClient.post<{ data: Group }>('/groups', payload);
      return { success: true, data: response.data.data };
    } catch (error: any) {
      return { success: false, error: error.message || 'Failed to create group' };
    }
  },

  updateGroup: async (id: string, payload: CreateAndUpdateGroupPayload): Promise<ApiResponse<Group>> => {
    try {
      const response = await httpClient.patch<{ data: Group }>(`/groups/${id}`, payload);
      return { success: true, data: response.data.data };
    } catch (error: any) {
      return { success: false, error: error.message || 'Failed to update group' };
    }
  },

  deleteGroup: async (id: string): Promise<ApiResponse<void>> => {
    try {
      await httpClient.delete(`/groups/${id}`);
      return { success: true };
    } catch (error: any) {
      return { success: false, error: error.message || 'Failed to delete group' };
    }
  },
};
