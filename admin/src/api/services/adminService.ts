import api from '@/api/client'
import type { ApiResponse, PaginatedResponse } from '@/types/common'
import type { AdminUser, AdminRole, SystemLog } from '@/types/admin'

export const getAdminProfile = () =>
    api.get<ApiResponse<AdminUser>>('/admin/me').then(r => r.data)

export const updateAdminProfile = (data: { email?: string; full_name?: string }) =>
    api.put<ApiResponse<AdminUser>>('/admin/me', data).then(r => r.data)

export interface AdminUserListParams {
    page?: number
    page_size?: number
    sort_by?: string
    sort_order?: string
    search?: string
    role_id?: string
    is_active?: boolean
}

export const listAdminUsers = (params: AdminUserListParams) =>
    api.get<ApiResponse<PaginatedResponse<AdminUser>>>('/admin/users', { params }).then(r => r.data)

export const getAdminUser = (id: string) =>
    api.get<ApiResponse<AdminUser>>(`/admin/users/${id}`).then(r => r.data)

export const createAdminUser = (data: {
    username: string
    email: string
    password: string
    full_name?: string
    role_id?: string
    is_active: boolean
}) =>
    api.post<ApiResponse<AdminUser>>('/admin/users', data).then(r => r.data)

export const updateAdminUser = (id: string, data: {
    email?: string
    full_name?: string
    role_id?: string
    is_active?: boolean
}) =>
    api.put<ApiResponse<AdminUser>>(`/admin/users/${id}`, data).then(r => r.data)

export const deleteAdminUser = (id: string) =>
    api.delete<ApiResponse<null>>(`/admin/users/${id}`).then(r => r.data)

export const deactivateAdminUser = (id: string) =>
    api.post<ApiResponse<AdminUser>>(`/admin/users/${id}/deactivate`).then(r => r.data)

export interface RoleListParams {
    page?: number
    page_size?: number
    sort_by?: string
    sort_order?: string
}

export const listAdminRoles = (params: RoleListParams) =>
    api.get<ApiResponse<PaginatedResponse<AdminRole>>>('/admin/roles', { params }).then(r => r.data)

export const getAdminRole = (id: string) =>
    api.get<ApiResponse<AdminRole>>(`/admin/roles/${id}`).then(r => r.data)

export const createAdminRole = (data: {
    role_name: string
    description: string
    permissions?: string
}) =>
    api.post<ApiResponse<AdminRole>>('/admin/roles', data).then(r => r.data)

export const updateAdminRole = (id: string, data: {
    description?: string
    permissions?: string
}) =>
    api.put<ApiResponse<AdminRole>>(`/admin/roles/${id}`, data).then(r => r.data)

export const deleteAdminRole = (id: string) =>
    api.delete<ApiResponse<null>>(`/admin/roles/${id}`).then(r => r.data)

export interface LogListParams {
    page?: number
    page_size?: number
    sort_by?: string
    sort_order?: string
    date_from?: string
    date_to?: string
    log_type?: string
    severity?: string
}

export const listSystemLogs = (params: LogListParams) =>
    api.get<ApiResponse<PaginatedResponse<SystemLog>>>('/admin/logs', { params }).then(r => r.data)

export const getSystemLog = (id: string) =>
    api.get<ApiResponse<SystemLog>>(`/admin/logs/${id}`).then(r => r.data)
