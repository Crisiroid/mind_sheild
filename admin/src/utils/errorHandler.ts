import { toast } from 'sonner'

export interface ApiError {
    message: string
    details?: string
    status?: number
}

export const extractErrorMessage = (error: any): string => {
    if (error.response?.data?.message) {
        return error.response.data.message
    }
    if (error.response?.data?.error) {
        return error.response.data.error
    }
    if (error.message) {
        return error.message
    }
    return 'خطای نامشخص رخ داده است'
}

export const handleApiError = (error: any, defaultMessage?: string): ApiError => {
    const message = extractErrorMessage(error)
    const status = error.response?.status

    if (status === 401) {
        return { message, status }
    }

    let displayMessage = message
    if (!message || message === 'Unknown Error') {
        switch (status) {
            case 400:
                displayMessage = 'درخواست نامعتبر است'
                break
            case 403:
                displayMessage = 'دسترسی مجاز نیست'
                break
            case 404:
                displayMessage = 'یافت نشد'
                break
            case 500:
                displayMessage = 'خطای سرور، لطفاً دوباره تلاش کنید'
                break
            default:
                displayMessage = defaultMessage || 'خطایی رخ داده است'
        }
    }

    toast.error(displayMessage, {
        description: error.response?.data?.details,
    })

    return {
        message: displayMessage,
        details: error.response?.data?.details,
        status,
    }
}

export const showSuccess = (message: string, description?: string) => {
    toast.success(message, { description })
}

export const showWarning = (message: string, description?: string) => {
    toast.warning(message, { description })
}

export const showInfo = (message: string, description?: string) => {
    toast.info(message, { description })
}
