export enum UserRole {
  USER = 'USER',
  ADMIN = 'ADMIN',
}

export interface User {
  id: string;
  email: string; // Primary identifier
  role: UserRole;
  avatar?: string;
  department?: string;
}
