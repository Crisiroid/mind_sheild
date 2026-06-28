export interface ApiResponse<T = any> {
    success: boolean
    message: string
    data: T
}

export interface PaginatedResponse<T = any> {
    data: T[]
    total: number
    page: number
    page_size: number
    pages: number
}

export interface PaginatedRequest {
    page?: number
    page_size?: number
    sort_by?: string
    sort_order?: 'asc' | 'desc'
    date_from?: string
    date_to?: string
    user_id?: string
    search?: string
}

export interface FilterRequest {
    [key: string]: any
}

export interface BaseEntity {
    id: string
    created_at: string
    updated_at: string
}
