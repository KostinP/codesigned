import { create } from 'zustand'

interface SessionUser {
  id: string
  username: string
  full_name: string
  role: string
}

interface SessionState {
  token: string | null
  user: SessionUser | null
  visitorId: string | null
  inactivityTimeout: number
  setToken: (token: string | null) => void
  setUser: (user: SessionUser | null) => void
  setVisitorId: (id: string) => void
  setInactivityTimeout: (seconds: number) => void
  reset: () => void
}

export const useSessionStore = create<SessionState>((set) => ({
  token: null,
  user: null,
  visitorId: null,
  inactivityTimeout: 3600,
  setToken: (token) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token || '')
    }
    set({ token })
  },
  setUser: (user) => set({ user }),
  setVisitorId: (visitorId) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('visitor_id', visitorId)
    }
    set({ visitorId })
  },
  setInactivityTimeout: (seconds) => set({ inactivityTimeout: seconds }),
  reset: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token')
    }
    set({ token: null, user: null })
  },
}))