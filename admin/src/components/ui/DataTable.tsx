import { TableLoading } from './LoadingSpinner'
import Pagination from './Pagination'

export interface Column<T = any> {
    key: string
    title: string
    render?: (item: T) => React.ReactNode
    sortable?: boolean
    width?: string
}

interface DataTableProps<T> {
    columns: Column<T>[]
    data: T[]
    loading?: boolean
    page?: number
    pageSize?: number
    total?: number
    pages?: number
    onPageChange?: (page: number) => void
    onPageSizeChange?: (size: number) => void
    onRowClick?: (item: T) => void
    emptyMessage?: string
    actions?: (item: T) => React.ReactNode
}

export default function DataTable<T extends { id?: string }>({
    columns,
    data,
    loading,
    page = 1,
    pageSize = 20,
    total = 0,
    pages = 1,
    onPageChange,
    onPageSizeChange,
    onRowClick,
    emptyMessage = 'داده‌ای یافت نشد',
    actions,
}: DataTableProps<T>) {
    if (loading) return <TableLoading />

    return (
        <div className="bg-white rounded-xl border border-neutral-200 shadow-sm overflow-hidden">
            <div className="overflow-x-auto">
                <table className="w-full min-w-[600px]">
                    <thead>
                        <tr className="border-b border-neutral-200 bg-neutral-50">
                            <th className="text-right px-4 md:px-6 py-3 md:py-4 font-semibold text-neutral-600 text-xs uppercase tracking-wider">#</th>
                            {columns.map(col => (
                                <th
                                    key={col.key}
                                    className={`text-right px-4 md:px-6 py-3 md:py-4 font-semibold text-neutral-600 text-xs uppercase tracking-wider ${col.width ? col.width : ''}`}
                                >
                                    {col.title}
                                </th>
                            ))}
                            {actions && <th className="text-right px-4 md:px-6 py-3 md:py-4 font-semibold text-neutral-600 text-xs uppercase tracking-wider">عملیات</th>}
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-neutral-100">
                        {data.length === 0 ? (
                            <tr>
                                <td colSpan={columns.length + 2} className="text-center py-12 md:py-16">
                                    <div className="flex flex-col items-center justify-center">
                                        <div className="w-12 h-12 md:w-16 md:h-16 mb-3 md:mb-4 rounded-full bg-neutral-100 flex items-center justify-center">
                                            <svg className="w-6 h-6 md:w-8 md:h-8 text-neutral-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                                            </svg>
                                        </div>
                                        <p className="text-neutral-600 font-medium text-sm md:text-base">{emptyMessage}</p>
                                    </div>
                                </td>
                            </tr>
                        ) : (
                            data.map((item, index) => (
                                <tr
                                    key={item.id || index}
                                    className="hover:bg-neutral-50 transition-colors duration-150 ${onRowClick ? 'cursor-pointer' : ''}"
                                    onClick={() => onRowClick?.(item)}
                                >
                                    <td className="px-4 md:px-6 py-3 md:py-4 text-neutral-500 text-xs md:text-sm font-medium">
                                        {(page - 1) * pageSize + index + 1}
                                    </td>
                                    {columns.map(col => (
                                        <td key={col.key} className="px-4 md:px-6 py-3 md:py-4 text-neutral-700 text-xs md:text-sm">
                                            {col.render ? col.render(item) : (item as any)[col.key] ?? '-'}
                                        </td>
                                    ))}
                                    {actions && (
                                        <td className="px-4 md:px-6 py-3 md:py-4" onClick={e => e.stopPropagation()}>
                                            {actions(item)}
                                        </td>
                                    )}
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
            {total > 0 && onPageChange && (
                <div className="border-t border-neutral-200 bg-neutral-50 px-4 md:px-6 py-3 md:py-4">
                    <Pagination
                        page={page}
                        pageSize={pageSize}
                        total={total}
                        pages={pages}
                        onPageChange={onPageChange}
                        onPageSizeChange={onPageSizeChange}
                    />
                </div>
            )}
        </div>
    )
}
