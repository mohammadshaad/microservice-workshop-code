export interface Task {
  id: string;
  title: string;
  description: string;
  userId: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface User {
  id: number;
  name: string;
  email: string;
} 