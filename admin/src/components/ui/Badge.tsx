import { HTMLAttributes } from 'react'
import { cn } from '@/utils/cn'

interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info'
    size?: 'sm' | 'md'
}

export function Badge({
    variant = 'default',
    size = 'md',
    className,
    children,
    ...props
}: BadgeProps) {
    const variants = {
        default: "bg-neutral-100 text-neutral-700 border border-neutral-200",
        success: "bg-success-50 text-success-700 border border-success-200",
        warning: "bg-warning-50 text-warning-700 border border-warning-200",
        error: "bg-error-50 text-error-700 border border-error-200",
        info: "bg-primary-50 text-primary-700 border border-primary-200",
    }

    const sizes = {
        sm: "px-2 py-0.5 text-xs",
        md: "px-2.5 py-1 text-sm",
    }

    return (
        <span
            className={cn(
                "inline-flex items-center font-semibold rounded-full",
                variants[variant],
                sizes[size],
                className
            )}
            {...props}
        >
            {children}
        </span>
    )
}
