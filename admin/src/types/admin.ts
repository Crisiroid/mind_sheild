export interface AdminUser {
    id: string
    username: string
    email: string
    full_name: string
    role_id: string
    is_active: boolean
    last_login: string
    created_at: string
    updated_at: string
}

export interface AdminLoginRequest {
    username: string
    password: string
}

export interface AdminLoginResponse {
    access_token: string
    refresh_token: string
    expires_in: number
    token_type: string
    admin_user: AdminUser
}

export interface AdminRefreshTokenRequest {
    refresh_token: string
}

export interface AdminRefreshTokenResponse {
    access_token: string
    refresh_token: string
    expires_in: number
    token_type: string
}

export interface AdminChangePasswordRequest {
    old_password: string
    new_password: string
}

export interface AdminCreateRequest {
    username: string
    email: string
    password: string
    full_name?: string
    role_id?: string
    is_active: boolean
}

export interface AdminUpdateRequest {
    email?: string
    full_name?: string
    role_id?: string
    is_active?: boolean
}

export interface AdminUpdateProfileRequest {
    email?: string
    full_name?: string
}

export interface AdminRole {
    id: string
    role_name: string
    description: string
    permissions: string
    created_at: string
    updated_at: string
}

export interface AdminRoleCreateRequest {
    role_name: string
    description: string
    permissions: string
}

export interface AdminRoleUpdateRequest {
    description?: string
    permissions?: string
}

export interface SystemLog {
    id: string
    log_type: string
    log_message: string
    user_id: string
    severity: string
    created_at: string
}
