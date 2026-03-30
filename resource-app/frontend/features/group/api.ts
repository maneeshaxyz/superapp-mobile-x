import { httpClient } from '../../api/client';
import { ApiResponse } from '../../api/types';
import { Group, CreateAndUpdateGroupPayload, GroupMember, AddUsersToGroupResult, RemoveUserFromGroupResult } from './types';

export const groupApi = {
  getGroups: async (): Promise<ApiResponse<Group[]>> => {
    try {
      const response = await httpClient.get<{ data: Group[] }>('/groups');
      return { success: true, data: response.data.data };
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Failed to fetch groups';
      return { success: false, error: message };
    }
  },

  createGroup: async (payload: CreateAndUpdateGroupPayload): Promise<ApiResponse<Group>> => {
    try {
      const response = await httpClient.post<{ data: Group }>('/groups', payload);
      return { success: true, data: response.data.data };
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Failed to create group';
      return { success: false, error: message };
    }
  },

  updateGroup: async (id: string, payload: CreateAndUpdateGroupPayload): Promise<ApiResponse<Group>> => {
    try {
      const response = await httpClient.patch<{ data: Group }>(`/groups/${id}`, payload);
      return { success: true, data: response.data.data };
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Failed to update group';
      return { success: false, error: message };
    }
  },

  deleteGroup: async (id: string): Promise<ApiResponse<void>> => {
    try {
      await httpClient.delete(`/groups/${id}`);
      return { success: true };
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Failed to delete group';
      return { success: false, error: message };
    }
  },

  // ── Membership ───────────────────────────────────────────────

  getGroupMembers: async (groupId: string): Promise<ApiResponse<GroupMember[]>> => {
    try {
      const response = await httpClient.get<{ data: GroupMember[] }>(`/groups/${groupId}/users`);
      return { success: true, data: response.data.data };
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Failed to fetch group members';
      return { success: false, error: message };
    }
  },

  addUsersToGroup: async (groupId: string, userIds: string[]): Promise<ApiResponse<AddUsersToGroupResult>> => {
    try {
      const response = await httpClient.post<{ data: AddUsersToGroupResult }>(`/groups/${groupId}/users`, { user_ids: userIds });
      return { success: true, data: response.data.data };
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Failed to add users to group';
      return { success: false, error: message };
    }
  },

  removeUserFromGroup: async (groupId: string, userId: string): Promise<ApiResponse<RemoveUserFromGroupResult>> => {
    try {
      const response = await httpClient.delete<{ data: RemoveUserFromGroupResult }>(`/groups/${groupId}/users/${userId}`);
      return { success: true, data: response.data.data };
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Failed to remove user from group';
      return { success: false, error: message };
    }
  },
};
