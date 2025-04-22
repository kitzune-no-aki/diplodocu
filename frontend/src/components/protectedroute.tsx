import { useAuth } from "../utils/AuthContext";
import { ReactNode } from "react";
import Unauthorized from "../components/Unauthorized"; // Your custom component

interface ProtectedRouteProps {
    children: ReactNode;
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
    const { isAuthenticated } = useAuth();

    if (!isAuthenticated) {
        return <Unauthorized />;
    }

    return <>{children}</>;
}