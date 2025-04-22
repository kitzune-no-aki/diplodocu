import { Link } from 'react-router-dom';
import { Lock, Home, LogIn } from 'lucide-react';
import { useAuth } from '../utils/AuthContext';

export default function UnauthorizedPage() {
    const { login } = useAuth();

    return (
        <div className="min-h-screen flex flex-col items-center justify-center bg-Dino-light p-4">
            <div className="w-full max-w-md bg-Dino-dark p-8 rounded-lg shadow-md text-center">
                <div className="mx-auto w-16 h-16 flex items-center justify-center bg-Warn-tomato border-2 border-Warn-fire rounded-full mb-4">
                    <Lock className="w-8 h-8 text-Warn-fire" />
                </div>

                <h1 className="text-2xl font-bold text-Oldrose mb-2">Access Denied</h1>
                <p className="text-Oldrose mb-6">
                    You don't have permission to view this page. Please log in with an authorized account.
                </p>

                <div className="flex flex-col space-y-3">
                    <button
                        onClick={login}
                        className="flex items-center justify-center gap-2 px-4 py-2 bg-Dino-light text-Aubergine rounded border-2 border-Dino-light hover:bg-Dino-dark hover:text-Oldrose transition"
                    >
                        <LogIn className="w-5 h-5" />
                        Sign In
                    </button>

                    <Link
                        to="/"
                        className="flex items-center justify-center gap-2 px-4 py-2 border border-Dino-light bg-Dino-light text-Aubergine rounded hover:bg-Dino-dark hover:text-Oldrose hover:border hover:border-Sand transition"
                    >
                        <Home className="w-5 h-5" />
                        Return Home
                    </Link>
                </div>
            </div>
        </div>
    );
}