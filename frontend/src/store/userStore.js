import { create } from 'zustand'
import { persist } from 'zustand/middleware'

const useUserStore = create(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      isAuthenticated: false,
      
      setToken: (token) => set({ token, isAuthenticated: !!token }),
      
      setUser: (user) => set({ user }),
      
      login: (token, user) => set({ 
        token, 
        user, 
        isAuthenticated: true 
      }),
      
      logout: () => set({ 
        token: null, 
        user: null, 
        isAuthenticated: false 
      }),
      
      updateUser: (userData) => {
        const { user } = get()
        set({ user: { ...user, ...userData } })
      },
    }),
    {
      name: 'user-storage',
    }
  )
)

export default useUserStore
