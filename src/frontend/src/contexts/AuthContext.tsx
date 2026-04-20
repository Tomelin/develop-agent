"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { UserProfile, userService } from "../services/userService";
import { authService } from "../services/authService";
import { useRouter } from "next/navigation";

interface AuthContextType {
  user: UserProfile | null;
  loading: boolean;
  login: (token: string, redirectUrl?: string) => Promise<void>;
  logout: () => Promise<void>;
  updateUser: (user: UserProfile) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem("@agency:token");
      if (!token) {
        setLoading(false);
        return;
      }

      try {
        const userData = await userService.getMe();
        setUser(userData);
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      } catch (error) {
        localStorage.removeItem("@agency:token");
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = async (token: string, redirectUrl: string = "/dashboard") => {
    localStorage.setItem("@agency:token", token);
    try {
      const userData = await userService.getMe();
      setUser(userData);
      router.push(redirectUrl);
    } catch (error) {
      console.error("Failed to fetch user data after login", error);
      localStorage.removeItem("@agency:token");
    }
  };

  const logout = async () => {
    try {
      await authService.logout();
    } catch (error) {
      console.error("Logout error", error);
    } finally {
      localStorage.removeItem("@agency:token");
      setUser(null);
      router.push("/login");
    }
  };

  const updateUser = (updatedUser: UserProfile) => {
    setUser(updatedUser);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, updateUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
