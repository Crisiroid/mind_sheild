import axios, { AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { useAuthStore } from '@/store/auth'
import { handleApiError } from '@/utils/errorHandler'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
})

let isRefreshing = false
let refreshSubscribers: Array<(token: string) => void> = []

const subscribeTokenRefresh = (callback: (token: string) => void) => {
    refreshSubscribers.push(callback)
}

const onRefreshed = (token: string) => {
    refreshSubscribers.forEach(callback => callback(token))
    refreshSubscribers = []
}

api.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        const { accessToken } = useAuthStore.getState()
        if (accessToken) {
            config.headers.Authorization = `Bearer ${accessToken}`
        }
        return config
    },
    (error: AxiosError) => {
        return Promise.reject(error)
    }
)

api.interceptors.response.use(
    (response: AxiosResponse) => {
        return response
    },
    async (error: AxiosError) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true

            const url = originalRequest.url || ''
            if (url.includes('/auth/')) {
                return Promise.reject(error)
            }

            try {
                const { refreshToken, setTokens, logout } = useAuthStore.getState()

                if (!refreshToken) {
                    logout()
                    if (window.location.pathname !== '/login') {
                        window.location.href = '/login'
                    }
                    return Promise.reject(error)
                }

                if (isRefreshing) {
                    return new Promise((resolve) => {
                        subscribeTokenRefresh((newToken: string) => {
                            originalRequest.headers.Authorization = `Bearer ${newToken}`
                            resolve(api(originalRequest))
                        })
                    })
                }

                isRefreshing = true

                const refreshResponse = await axios.post(
                    `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/auth/admin/refresh`,
                    { refresh_token: refreshToken }
                )

                const { access_token, refresh_token } = refreshResponse.data.data

                setTokens(access_token, refresh_token)

                onRefreshed(access_token)

                isRefreshing = false

                originalRequest.headers.Authorization = `Bearer ${access_token}`
                return api(originalRequest)
            } catch (refreshError) {
                isRefreshing = false
                refreshSubscribers = []
                const { logout } = useAuthStore.getState()
                logout()
                if (window.location.pathname !== '/login') {
                    window.location.href = '/login'
                }
                return Promise.reject(refreshError)
            }
        }

        handleApiError(error)
        return Promise.reject(error)
    }
)

export default api
