import { HTMLAttributes } from 'react'
import { cn } from '@/utils/cn'

interface CardProps extends HTMLAttributes<HTMLDivElement> {
    variant?: 'default' | 'elevated' | 'bordered'
    hover?: boolean
}

export function Card({
    variant = 'default',
    hover = false,
    className,
    children,
    ...props
}: CardProps) {
    const variants = {
        default: "bg-white border border-neutral-200 shadow-sm",
        elevated: "bg-white shadow-soft-lg border border-neutral-100",
        bordered: "bg-white border-2 border-neutral-200",
    }

    return (
        <div
            className={cn(
                "rounded-xl p-6 transition-all duration-200",
                variants[variant],
                hover && "hover:shadow-lg hover:-translate-y-1 hover:border-neutral-300",
                className
            )}
            {...props}
        >
            {children}
        </div>
    )
}

export function CardHeader({
    className,
    children,
    ...props
}: HTMLAttributes<HTMLDivElement>) {
    return (
        <div className={cn("mb-5 pb-4 border-b border-neutral-100", className)} {...props}>
            {children}
        </div>
    )
}

export function CardTitle({
    className,
    children,
    ...props
}: HTMLAttributes<HTMLHeadingElement>) {
    return (
        <h3 className={cn("text-lg font-bold tracking-tight text-neutral-900", className)} {...props}>
            {children}
        </h3>
    )
}

export function CardDescription({
    className,
    children,
    ...props
}: HTMLAttributes<HTMLParagraphElement>) {
    return (
        <p className={cn("text-sm text-neutral-500 mt-1", className)} {...props}>
            {children}
        </p>
    )
}

export function CardContent({
    className,
    children,
    ...props
}: HTMLAttributes<HTMLDivElement>) {
    return (
        <div className={cn("pt-2", className)} {...props}>
            {children}
        </div>
    )
}
