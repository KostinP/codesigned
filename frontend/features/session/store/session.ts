import { create } from 'zustand'

interface SessionState {
  token: string | null
  user: any | null
  visitorId: string | null
  inactivityTimeout: number
  setToken: (token: string | null) => void
  setUser: (user: any | null) => void
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
    localStorage.setItem('token', token || '')
    set({ token })
  },
  setUser: (user) => set({ user }),
  setVisitorId: (visitorId) => {
    localStorage.setItem('visitor_id', visitorId)
    set({ visitorId })
  },
  setInactivityTimeout: (seconds) => set({ inactivityTimeout: seconds }),
  reset: () => {
    localStorage.removeItem('token')
    set({ token: null, user: null })
  },
}))
