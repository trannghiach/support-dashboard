import { jwtDecode } from "jwt-decode";

const TOKEN_KEY = "support_dashboard_token";

type TokenPayload = {
    user_id: number;
    role: "customer" | "agent" | "admin";
    exp: number;
    iat: number;
};

export function saveToken(token: string) {
    if (typeof window === "undefined") return;
    localStorage.setItem(TOKEN_KEY, token);
}

export function getToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem(TOKEN_KEY);
}

export function removeToken() {
    if (typeof window === "undefined") return;
    localStorage.removeItem(TOKEN_KEY);
}

export function getAuthPayload(): TokenPayload | null {
    const token = getToken();
    if (!token) return null;

    try {
        return jwtDecode<TokenPayload>(token);
    } catch {
        return null;
    }
}