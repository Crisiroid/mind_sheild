import { InputHTMLAttributes, forwardRef } from 'react'
import { cn } from '@/utils/cn'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    label?: string
    error?: string
    icon?: React.ReactNode
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
    ({ label, error, icon, className, id, ...props }, ref) => {
        const inputId = id || label?.toLowerCase().replace(/\s+/g, '-')

        return (
            <div className="space-y-2">
                {label && (
                    <label htmlFor={inputId} className="block text-sm font-semibold text-neutral-700 mb-1.5">
                        {label}
                    </label>
                )}
                <div className="relative">
                    {icon && (
                        <div className="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 pointer-events-none">
                            {icon}
                        </div>
                    )}
                    <input
                        ref={ref}
                        id={inputId}
                        className={cn(
                            "w-full px-4 py-2.5 text-sm",
                            "bg-white border border-neutral-300 rounded-lg",
                            "placeholder:text-neutral-400",
                            "transition-all duration-200 ease-out",
                            "focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500",
                            "disabled:bg-neutral-50 disabled:cursor-not-allowed",
                            icon && "pr-10",
                            error && "border-error-500 focus:ring-error-500/20 focus:border-error-500",
                            className
                        )}
                        {...props}
                    />
                </div>
                {error && (
                    <p className="text-sm text-error-600 flex items-center gap-1.5 mt-1">
                        <svg className="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                            <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                        </svg>
                        {error}
                    </p>
                )}
            </div>
        )
    }
)

Input.displayName = 'Input'
