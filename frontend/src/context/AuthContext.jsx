import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { api } from "../api/client";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [status, setStatus] = useState("loading"); // loading | guest | authed

  // session check on mount (reading the httpOnly cookie via /auth/me);
  // there is no external store to subscribe to instead, so the effect's
  // job IS to fetch and set the initial auth state.
  useEffect(() => {
    const controller = new AbortController();

    api
      .me()
      .then((u) => {
        if (controller.signal.aborted) return;
        setUser(u);
        setStatus("authed");
      })
      .catch(() => {
        if (controller.signal.aborted) return;
        setStatus("guest");
      });

    return () => controller.abort();
  }, []);

  const register = useCallback(async (payload) => {
    const res = await api.register(payload);
    setUser(res.user);
    setStatus("authed");
    return res.user;
  }, []);

  const login = useCallback(async (payload) => {
    const res = await api.login(payload);
    setUser(res.user);
    setStatus("authed");
    return res.user;
  }, []);

  const logout = useCallback(async () => {
    try {
      await api.logout();
    } finally {
      setUser(null);
      setStatus("guest");
    }
  }, []);

  return (
    <AuthContext.Provider value={{ user, status, register, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// useAuth is exported alongside the AuthProvider component (standard
// context+hook co-location); splitting it into a second file for one hook
// would add indirection without benefit for a project this size.
// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
