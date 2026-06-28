interface PaginationProps {
    page: number
    pageSize: number
    total: number
    pages: number
    onPageChange: (page: number) => void
    onPageSizeChange?: (size: number) => void
}

export default function Pagination({ page, pageSize, total, pages, onPageChange, onPageSizeChange }: PaginationProps) {
    const startItem = (page - 1) * pageSize + 1
    const endItem = Math.min(page * pageSize, total)

    const getPageNumbers = () => {
        const delta = 2
        const range: (number | string)[] = []
        for (let i = 1; i <= pages; i++) {
            if (i === 1 || i === pages || (i >= page - delta && i <= page + delta)) {
                range.push(i)
            } else if (range[range.length - 1] !== '...') {
                range.push('...')
            }
        }
        return range
    }

    return (
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            <div className="text-sm text-neutral-600 font-medium">
                نمایش {startItem} تا {endItem} از {total.toLocaleString()} مورد
            </div>
            <div className="flex items-center gap-3">
                {onPageSizeChange && (
                    <select
                        value={pageSize}
                        onChange={e => onPageSizeChange(Number(e.target.value))}
                        className="text-sm border border-neutral-300 rounded-lg px-3 py-2 bg-white hover:border-neutral-400 transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    >
                        {[20, 50, 100].map(size => (
                            <option key={size} value={size}>{size} در صفحه</option>
                        ))}
                    </select>
                )}
                <div className="flex items-center gap-2">
                    <button
                        onClick={() => onPageChange(page - 1)}
                        disabled={page <= 1}
                        className="px-4 py-2 text-sm font-medium border border-neutral-300 rounded-lg bg-white text-neutral-700 hover:bg-neutral-50 hover:border-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white disabled:hover:border-neutral-300 transition-all duration-150"
                    >
                        قبلی
                    </button>
                    {getPageNumbers().map((p, i) =>
                        typeof p === 'string' ? (
                            <span key={`dots-${i}`} className="px-2 py-2 text-sm text-neutral-400">...</span>
                        ) : (
                            <button
                                key={p}
                                onClick={() => onPageChange(p)}
                                className={`min-w-[40px] px-3 py-2 text-sm font-medium border rounded-lg transition-all duration-150 ${p === page
                                        ? 'bg-primary-600 text-white border-primary-600 shadow-sm'
                                        : 'border-neutral-300 bg-white text-neutral-700 hover:bg-neutral-50 hover:border-neutral-400'
                                    }`}
                            >
                                {p}
                            </button>
                        )
                    )}
                    <button
                        onClick={() => onPageChange(page + 1)}
                        disabled={page >= pages}
                        className="px-4 py-2 text-sm font-medium border border-neutral-300 rounded-lg bg-white text-neutral-700 hover:bg-neutral-50 hover:border-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white disabled:hover:border-neutral-300 transition-all duration-150"
                    >
                        بعدی
                    </button>
                </div>
            </div>
        </div>
    )
}
