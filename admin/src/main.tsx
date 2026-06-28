import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { Toaster } from 'sonner'
import App from './App'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <BrowserRouter>
            <App />
            <Toaster
                position="top-right"
                richColors
                closeButton
                toastOptions={{
                    className: 'font-sans',
                    style: {
                        fontFamily: 'Vazirmatn, system-ui, sans-serif',
                        direction: 'rtl',
                    },
                }}
            />
        </BrowserRouter>
    </React.StrictMode>,
)
