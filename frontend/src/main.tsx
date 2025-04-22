import { createRoot } from 'react-dom/client'
import React from "react";
import './index.css'
import App from './App.tsx'
import { AuthProvider } from './utils/AuthContext.tsx'

createRoot(document.getElementById('root')!).render(
    <AuthProvider>
        <React.StrictMode>
            <App />
        </React.StrictMode>
    </AuthProvider>,
)
