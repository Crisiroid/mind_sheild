import api from '@/api/client'
import type { ApiResponse } from '@/types/common'

export const getDashboard = () =>
    api.get<ApiResponse<any>>('/reports/dashboard').then(r => r.data)

export const getUserActivity = (params?: { date_from?: string; date_to?: string }) =>
    api.get<ApiResponse<any>>('/reports/user-activity', { params }).then(r => r.data)

export const getStressAnalytics = (params?: { date_from?: string; date_to?: string }) =>
    api.get<ApiResponse<any>>('/reports/stress-analytics', { params }).then(r => r.data)

export const getBodyTensionReport = (params?: { date_from?: string; date_to?: string }) =>
    api.get<ApiResponse<any>>('/reports/body-tension', { params }).then(r => r.data)

export const getCognitivePatterns = (params?: { date_from?: string; date_to?: string }) =>
    api.get<ApiResponse<any>>('/reports/cognitive-patterns', { params }).then(r => r.data)

export const getMoodTrends = (params?: { date_from?: string; date_to?: string }) =>
    api.get<ApiResponse<any>>('/reports/mood-trends', { params }).then(r => r.data)

export const getEngagement = (params?: { date_from?: string; date_to?: string }) =>
    api.get<ApiResponse<any>>('/reports/engagement', { params }).then(r => r.data)

export const getWeeklyProgress = (params?: { date_from?: string; date_to?: string }) =>
    api.get<ApiResponse<any>>('/reports/weekly-progress', { params }).then(r => r.data)

export const exportData = (params?: { date_from?: string; date_to?: string; user_id?: string }) =>
    api.get<ApiResponse<any>>('/reports/export', { params }).then(r => r.data)

export const listUserReports = (params?: { page?: number; page_size?: number }) =>
    api.get<ApiResponse<any>>('/admin/reports', { params }).then(r => r.data)

export const createUserReport = (data: {
    report_type: string
    report_date: string
    total_users: number
    active_users: number
    crisis_alerts_count: number
}) =>
    api.post<ApiResponse<any>>('/admin/reports', data).then(r => r.data)

export const deleteUserReport = (id: string) =>
    api.delete<ApiResponse<null>>(`/admin/reports/${id}`).then(r => r.data)

export interface MediaListParams {
    page?: number
    page_size?: number
    week_number?: number
    file_type?: string
    is_active?: boolean
}

export const listMediaContent = (params: MediaListParams) =>
    api.get<ApiResponse<any>>('/media/weekly', { params }).then(r => r.data)

export const getMediaContent = (id: string) =>
    api.get<ApiResponse<any>>(`/media/weekly/${id}`).then(r => r.data)

export const uploadMediaContent = (formData: FormData) =>
    api.post<ApiResponse<any>>('/media/weekly', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
    }).then(r => r.data)

export const updateMediaContent = (id: string, data: any) =>
    api.put<ApiResponse<any>>(`/media/weekly/${id}`, data).then(r => r.data)

export const deleteMediaContent = (id: string) =>
    api.delete<ApiResponse<null>>(`/media/weekly/${id}`).then(r => r.data)

export const getUserActivityTrend = (params: { date_from: string; date_to: string }) =>
    api.get<ApiResponse<any>>('/users/activity-trend', { params }).then(r => r.data)

export const getUserLoginAnalytics = (params: { date_from: string; date_to: string }) =>
    api.get<ApiResponse<any>>('/users/login-analytics', { params }).then(r => r.data)

export const getUserEngagement = (params: { date_from: string; date_to: string }) =>
    api.get<ApiResponse<any>>('/users/engagement', { params }).then(r => r.data)

export const getInactiveUsers = (days?: number) =>
    api.get<ApiResponse<any>>('/users/inactive', { params: { days: days || 30 } }).then(r => r.data)

export const changeAdminPassword = (data: { old_password: string; new_password: string }) =>
    api.post<ApiResponse<null>>('/auth/admin/change-password', data).then(r => r.data)
