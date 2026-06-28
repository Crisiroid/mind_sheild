export interface User {
    id: string
    phone_number: string
    registration_date: string
    last_login: string
    login_count: number
    agreement_accepted: boolean
    agreement_accepted_at: string
    cloud_sync_enabled: boolean
    do_not_disturb_enabled: boolean
    dnd_start_time: string
    dnd_end_time: string
    android_version: string
    app_version: string
    created_at: string
    updated_at: string
}

export interface UserCreateRequest {
    phone_number: string
    android_version?: string
    app_version?: string
}

export interface UserUpdateRequest {
    phone_number?: string
    cloud_sync_enabled?: boolean
    do_not_disturb_enabled?: boolean
    dnd_start_time?: string
    dnd_end_time?: string
}
