import api from '@/api/client'
import type { ApiResponse, PaginatedResponse } from '@/types/common'
import type { User } from '@/types/user'

export interface UserListParams {
    page?: number
    page_size?: number
    sort_by?: string
    sort_order?: string
    date_from?: string
    date_to?: string
    search?: string
    agreement_accepted?: boolean
}

export const listUsers = (params: UserListParams) =>
    api.get<ApiResponse<PaginatedResponse<User>>>('/users', { params }).then(r => r.data)

export const getUser = (id: string) =>
    api.get<ApiResponse<User>>(`/users/${id}`).then(r => r.data)

export const createUser = (data: { phone_number: string; android_version?: string; app_version?: string }) =>
    api.post<ApiResponse<User>>('/users', data).then(r => r.data)

export const updateUser = (id: string, data: {
    phone_number?: string
    cloud_sync_enabled?: boolean
    do_not_disturb_enabled?: boolean
    dnd_start_time?: string
    dnd_end_time?: string
}) =>
    api.put<ApiResponse<User>>(`/users/${id}`, data).then(r => r.data)

export const deleteUser = (id: string) =>
    api.delete<ApiResponse<null>>(`/users/${id}`).then(r => r.data)

export const getUserByPhone = (phone: string) =>
    api.get<ApiResponse<User>>('/users/by-phone', { params: { phone_number: phone } }).then(r => r.data)

export const acceptAgreement = (id: string) =>
    api.post<ApiResponse<User>>(`/users/${id}/accept-agreement`).then(r => r.data)

export interface UserStats {
    total_users: number
    active_users: number
    new_users_today: number
    new_users_this_week: number
    new_users_this_month: number
    agreement_rate: number
    avg_login_count: number
}

export const getUserStats = () =>
    api.get<ApiResponse<UserStats>>('/users/stats').then(r => r.data)

export const getAgreementStats = () =>
    api.get<ApiResponse<{ total_users: number; agreed_users: number; agreement_rate: number }>>('/users/agreement-stats').then(r => r.data)

export const getAppVersionDistribution = () =>
    api.get<ApiResponse<any>>('/users/app-version-distribution').then(r => r.data)

export const exportUsers = (params?: { date_from?: string; date_to?: string; user_id?: string }) =>
    api.get<ApiResponse<any>>('/users/export', { params }).then(r => r.data)
