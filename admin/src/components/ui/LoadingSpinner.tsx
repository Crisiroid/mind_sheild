export function LoadingSpinner({ size = 'md', className = '' }: { size?: 'sm' | 'md' | 'lg'; className?: string }) {
    const sizes = { sm: 'h-5 w-5', md: 'h-8 w-8', lg: 'h-12 w-12' }
    return (
        <div className={`flex justify-center items-center ${className}`}>
            <div className={`${sizes[size]} animate-spin rounded-full border-4 border-primary-200 border-t-primary-600 shadow-sm`} />
        </div>
    )
}

export function PageLoading() {
    return (
        <div className="flex flex-col justify-center items-center py-20 gap-4">
            <div className="relative">
                <LoadingSpinner size="lg" />
                <div className="absolute inset-0 animate-pulse-soft">
                    <div className="h-12 w-12 rounded-full bg-primary-500/10" />
                </div>
            </div>
            <p className="text-sm text-neutral-500 font-medium">در حال بارگذاری...</p>
        </div>
    )
}

export function TableLoading() {
    return (
        <div className="bg-white rounded-xl border border-neutral-200 shadow-sm">
            <div className="flex justify-center items-center py-20">
                <div className="flex flex-col items-center gap-4">
                    <div className="relative">
                        <LoadingSpinner size="md" />
                        <div className="absolute inset-0 animate-pulse-soft">
                            <div className="h-8 w-8 rounded-full bg-primary-500/10" />
                        </div>
                    </div>
                    <p className="text-sm text-neutral-500 font-medium">در حال بارگذاری...</p>
                </div>
            </div>
        </div>
    )
}

export function Skeleton({ className = '' }: { className?: string }) {
    return (
        <div className={`skeleton rounded-lg ${className}`} />
    )
}
