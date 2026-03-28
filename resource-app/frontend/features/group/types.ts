export interface Group {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAndUpdateGroupPayload {
  name: string;
  description: string;
}
