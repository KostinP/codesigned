import { create } from 'zustand'

interface AuthUser {
  id: string
  username: string
  full_name: string
  role: string
}

interface AuthState {
  user: AuthUser | null
  setUser: (user: AuthUser | null) => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  setUser: (user) => set({ user }),
}))