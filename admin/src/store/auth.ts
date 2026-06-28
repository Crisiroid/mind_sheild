import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { AdminUser, AdminLoginRequest } from '@/types/admin'
import api from '@/api/client'
import { showSuccess } from '@/utils/errorHandler'

interface AuthState {
    accessToken: string | null
    refreshToken: string | null
    adminUser: AdminUser | null
    isAuthenticated: boolean
    isLoading: boolean

    login: (credentials: AdminLoginRequest) => Promise<void>
    logout: () => void
    setTokens: (access: string, refresh: string) => void
    updateProfile: (profile: Partial<AdminUser>) => void
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            accessToken: null,
            refreshToken: null,
            adminUser: null,
            isAuthenticated: false,
            isLoading: false,

            login: async (credentials: AdminLoginRequest) => {
                set({ isLoading: true })
                try {
                    const response = await api.post('/auth/admin/login', credentials)
                    const { access_token, refresh_token, admin_user } = response.data.data

                    set({
                        accessToken: access_token,
                        refreshToken: refresh_token,
                        adminUser: admin_user,
                        isAuthenticated: true,
                        isLoading: false,
                    })

                    showSuccess('ورود موفقیت‌آمیز')
                } catch (error) {
                    set({ isLoading: false })
                    throw error
                }
            },

            logout: () => {
                set({
                    accessToken: null,
                    refreshToken: null,
                    adminUser: null,
                    isAuthenticated: false,
                })
                showSuccess('خروج موفقیت‌آمیز')
            },

            setTokens: (access: string, refresh: string) => {
                set({
                    accessToken: access,
                    refreshToken: refresh,
                })
            },

            updateProfile: (profile: Partial<AdminUser>) => {
                set((state) => ({
                    adminUser: state.adminUser ? { ...state.adminUser, ...profile } : null,
                }))
            },
        }),
        {
            name: 'auth-storage',
            partialize: (state) => ({
                accessToken: state.accessToken,
                refreshToken: state.refreshToken,
                adminUser: state.adminUser,
                isAuthenticated: state.isAuthenticated,
            }),
        }
    )
)
