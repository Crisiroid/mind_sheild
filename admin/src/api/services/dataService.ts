import api from '@/api/client'
import type { ApiResponse, PaginatedResponse } from '@/types/common'

export interface ListParams {
    page?: number
    page_size?: number
    sort_by?: string
    sort_order?: string
    date_from?: string
    date_to?: string
    user_id?: string
    search?: string
    [key: string]: any
}

export const listCalendars = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/calendars', { params }).then(r => r.data)

export const listEmotionInteractions = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/emotion-interactions', { params }).then(r => r.data)

export const listStressEvents = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/stress-events', { params }).then(r => r.data)

export const listBodyTensionMaps = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/body-tension-maps', { params }).then(r => r.data)

export const listBreathingSessions = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/breathing-sessions', { params }).then(r => r.data)

export const listCognitiveGames = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/cognitive-games', { params }).then(r => r.data)

export const listMentalMusts = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/mental-musts', { params }).then(r => r.data)

export const listNegativeThoughts = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/negative-thoughts', { params }).then(r => r.data)

export const listMindCourtEvidence = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/mind-court-evidence', { params }).then(r => r.data)

export const listConflictExercises = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/conflict-exercises', { params }).then(r => r.data)

export const listMoodTracker = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/mood-tracker', { params }).then(r => r.data)

export const listRolesValues = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/roles-values', { params }).then(r => r.data)

export const listSkyThoughts = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/sky-thoughts', { params }).then(r => r.data)

export const listMindfulTimers = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/mindful-timers', { params }).then(r => r.data)

export const listAcceptanceExercises = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/acceptance-exercises', { params }).then(r => r.data)

export const listWeeklyReports = (params: ListParams) =>
    api.get<ApiResponse<PaginatedResponse<any>>>('/reports/weekly', { params }).then(r => r.data)
