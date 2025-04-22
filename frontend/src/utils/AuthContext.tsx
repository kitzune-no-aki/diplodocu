import { createContext, useContext, useEffect, useState, ReactNode } from "react";
import keycloak from "../utils/Keycloak";

interface UserProfile {
    username?: string;
    email?: string;
    firstName?: string;
    lastName?: string;
}

interface AuthState {
    isAuthenticated: boolean;
    isLoading: boolean;
    userProfile: UserProfile | null;
    roles: string[];
}

interface AuthContextType extends AuthState {
    login: () => void;
    logout: () => void;
}

// Create context with initial value
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Props type for AuthProvider
interface AuthProviderProps {
    children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
    const [authState, setAuthState] = useState<AuthState>({
        isAuthenticated: false,
        isLoading: true,
        userProfile: null,
        roles: [],
    });

    useEffect(() => {
        const initializeAuth = async () => {
            try {
                const authenticated = await keycloak.init({
                    onLoad: "check-sso",
                    checkLoginIframe: false,
                });

                if (authenticated) {
                    const userProfile = await keycloak.loadUserProfile();

                    setAuthState({
                        isAuthenticated: true,
                        isLoading: false,
                        userProfile,
                        roles: keycloak.tokenParsed?.realm_access?.roles || [],
                    });
                } else {
                    setAuthState((prev) => ({ ...prev, isLoading: false }));
                }
            } catch (error) {
                console.error("Auth initialization failed:", error);
                setAuthState((prev) => ({ ...prev, isLoading: false }));
            }
        };

        initializeAuth();

        keycloak.onTokenExpired = () => {
            keycloak.updateToken(30).catch(() => {
                keycloak.logout();
            });
        };

        return () => {
            // @ts-ignore
            keycloak.onTokenExpired = null;
        };
    }, []);

    const value: AuthContextType = {
        ...authState,
        login: () =>
            keycloak.login({ redirectUri: window.location.origin + "/mytable" }),
        logout: () => {
            keycloak.logout({ redirectUri: window.location.origin });
        },
    };

    return (
        <AuthContext.Provider value={value}>
            {!authState.isLoading && children}
        </AuthContext.Provider>
    );
}

export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};