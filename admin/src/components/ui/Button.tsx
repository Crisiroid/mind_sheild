import { ButtonHTMLAttributes } from 'react'
import { cn } from '@/utils/cn'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger'
    size?: 'sm' | 'md' | 'lg'
    loading?: boolean
}

export function Button({
    variant = 'primary',
    size = 'md',
    loading = false,
    className,
    children,
    disabled,
    ...props
}: ButtonProps) {
    const baseStyles = "inline-flex items-center justify-center font-semibold rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed active:scale-95"

    const variants = {
        primary: "text-white focus:ring-primary-500 shadow-md hover:shadow-lg btn-lift" +
            " bg-gradient-to-br from-primary-500 to-primary-600 hover:from-primary-600 hover:to-primary-700",
        secondary: "bg-neutral-100 text-neutral-700 hover:bg-neutral-200 focus:ring-neutral-500 border border-neutral-200",
        outline: "border-2 border-primary-500 text-primary-600 hover:bg-primary-50 focus:ring-primary-500 bg-transparent",
        ghost: "text-neutral-700 hover:bg-neutral-100 focus:ring-neutral-500 bg-transparent",
        danger: "text-white focus:ring-error-500 shadow-sm hover:shadow-md btn-lift" +
            " bg-gradient-to-br from-error-500 to-error-600 hover:from-error-600 hover:to-error-700",
    }

    const sizes = {
        sm: "px-3 py-1.5 text-sm",
        md: "px-4 py-2.5 text-sm",
        lg: "px-6 py-3 text-base",
    }

    return (
        <button
            className={cn(baseStyles, variants[variant], sizes[size], className)}
            disabled={disabled || loading}
            {...props}
        >
            {loading && (
                <svg className="animate-spin ml-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
            )}
            <span className={loading ? 'opacity-90' : ''}>{children}</span>
        </button>
    )
}
