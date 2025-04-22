import { Link } from 'react-router-dom';
import { Rabbit, SquareLibrary } from 'lucide-react';
import { useAuth } from "../utils/AuthContext";

export default function Navbar() {
    const { isAuthenticated, login, logout } = useAuth();

    return (
        <div className="fixed bottom-0 w-full bg-Dino-dark shadow-md border-t border-Dino-dark z-50">
            <ul className="flex justify-around py-2">
                <li className="text-Oldrose hover:text-Aubergine">
                    <Link to='/' className="flex flex-col items-center">
                        <Rabbit className="w-6 h-6" />
                        <span className="text-xs">Home</span>
                    </Link>
                </li>
                <li className="text-Oldrose hover:text-Aubergine">
                    <Link to='/mytable' className="flex flex-col items-center">
                        <SquareLibrary className="w-6 h-6" />
                        <span className="text-xs">My Table</span>
                    </Link>
                </li>
                <li className="flex flex-col items-center">
                    {isAuthenticated ? (
                        <button
                            onClick={logout}
                            className="px-4 py-2 text-Oldrose hover:text-Aubergine font-medium transition-colors"
                        >
                            Logout
                        </button>
                    ) : (
                        <button
                            onClick={login}
                            className="px-4 py-2 text-Oldrose hover:text-Aubergine font-medium transition-colors"
                        >
                            Login
                        </button>
                    )}
                </li>
            </ul>
        </div>
    )
}