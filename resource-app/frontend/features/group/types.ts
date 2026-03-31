export interface Group {
  id: string;
  name: string;
  description: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateGroupPayload {
  name: string;
  description: string;
  userIds?: string[];
}

export interface UpdateGroupPayload {
  name: string;
  description: string;
}

export interface GroupMember {
  id: string;
  name: string;
  email: string;
}

export interface AddUsersToGroupResult {
  groupId: string;
  addedUsers: { userId: string }[];
}

export interface RemoveUserFromGroupResult {
  groupId: string;
  userId: string;
}
