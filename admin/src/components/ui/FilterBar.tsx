interface FilterBarProps {
    userId?: string
    onUserIdChange?: (id: string) => void
    dateFrom?: string
    dateTo?: string
    onDateFromChange?: (date: string) => void
    onDateToChange?: (date: string) => void
    search?: string
    onSearchChange?: (search: string) => void
    extraFilters?: React.ReactNode
    onReset?: () => void
}

export default function FilterBar({
    userId,
    onUserIdChange,
    dateFrom,
    dateTo,
    onDateFromChange,
    onDateToChange,
    search,
    onSearchChange,
    extraFilters,
    onReset,
}: FilterBarProps) {
    return (
        <div className="bg-white rounded-lg shadow p-3 md:p-4 mb-4 md:mb-6">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
                {onSearchChange && (
                    <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">جستجو</label>
                        <input
                            type="text"
                            value={search || ''}
                            onChange={e => onSearchChange(e.target.value)}
                            placeholder="جستجو..."
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                    </div>
                )}
                {onUserIdChange && (
                    <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">شناسه کاربر</label>
                        <input
                            type="text"
                            value={userId || ''}
                            onChange={e => onUserIdChange(e.target.value)}
                            placeholder="شناسه کاربر..."
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                    </div>
                )}
                {onDateFromChange && (
                    <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">از تاریخ</label>
                        <input
                            type="date"
                            value={dateFrom || ''}
                            onChange={e => onDateFromChange(e.target.value)}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                    </div>
                )}
                {onDateToChange && (
                    <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">تا تاریخ</label>
                        <input
                            type="date"
                            value={dateTo || ''}
                            onChange={e => onDateToChange(e.target.value)}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                    </div>
                )}
                {extraFilters}
            </div>
            {onReset && (
                <div className="mt-3 md:mt-4 flex justify-end">
                    <button
                        onClick={onReset}
                        className="text-sm text-gray-600 hover:text-gray-900 px-3 py-1.5 md:px-3 md:py-2 border border-gray-300 rounded-md hover:bg-gray-50 w-full sm:w-auto"
                    >
                        پاک کردن فیلترها
                    </button>
                </div>
            )}
        </div>
    )
}
