import { useState } from 'react'
import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import Header from './Header'

export default function MainLayout({ children }: { children?: React.ReactNode }) {
    const [sidebarOpen, setSidebarOpen] = useState(false)

    return (
        <div className="flex h-screen overflow-hidden" style={{ background: '#f7fafc' }}>
            {/* Mobile sidebar overlay */}
            {sidebarOpen && (
                <div
                    className="fixed inset-0 bg-black/50 z-40 lg:hidden"
                    onClick={() => setSidebarOpen(false)}
                />
            )}

            {/* Sidebar */}
            <div className={`
                fixed lg:static inset-y-0 right-0 z-50 transform transition-transform duration-300 ease-in-out
                ${sidebarOpen ? 'translate-x-0' : 'translate-x-full lg:translate-x-0'}
            `}>
                <Sidebar onClose={() => setSidebarOpen(false)} />
            </div>

            <div className="flex-1 flex flex-col overflow-hidden">
                <Header onMenuToggle={() => setSidebarOpen(!sidebarOpen)} />
                <main className="flex-1 overflow-y-auto p-4 md:p-6 lg:p-8">
                    {children || <Outlet />}
                </main>
            </div>
        </div>
    )
}
