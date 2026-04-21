"use client";

import { useAuth } from "@/contexts/AuthContext";
import { PrivateRoute } from "@/components/auth/PrivateRoute";
import Link from "next/link";
import { LayoutDashboard, User as UserIcon, LogOut, Bot, WandSparkles, ReceiptText, ShieldCheck, BarChart3, Building2, Store, Map, CreditCard } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { NotificationBell } from "@/components/dashboard/NotificationBell";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, logout } = useAuth();

  return (
    <PrivateRoute>
      <div className="min-h-screen bg-background flex flex-col">
        {/* Header/Nav */}
        <header className="sticky top-0 z-40 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div className="container mx-auto flex h-16 items-center justify-between px-4">
            <div className="flex items-center gap-6">
              <Link href="/dashboard" className="flex items-center space-x-2">
                <span className="font-bold text-primary sm:inline-block">
                  Agency AI
                </span>
              </Link>
              <nav className="hidden md:flex gap-6">
                <Link
                  href="/dashboard"
                  className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                >
                  <LayoutDashboard className="mr-2 h-4 w-4" />
                  Dashboard
                </Link>
                <Link
                  href="/dashboard/agents"
                  className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                >
                  <Bot className="mr-2 h-4 w-4" />
                  Agentes
                </Link>
                <Link
                  href="/dashboard/prompts"
                  className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                >
                  <WandSparkles className="mr-2 h-4 w-4" />
                  Prompts
                </Link>
                <Link
                  href="/dashboard/billing"
                  className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                >
                  <ReceiptText className="mr-2 h-4 w-4" />
                  Billing
                </Link>
                <Link
                  href="/dashboard/organization"
                  className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                >
                  <Building2 className="mr-2 h-4 w-4" />
                  Organização
                </Link>
                <Link
                  href="/dashboard/marketplace"
                  className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                >
                  <Store className="mr-2 h-4 w-4" />
                  Marketplace
                </Link>
                <Link
                  href="/dashboard/roadmap"
                  className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                >
                  <Map className="mr-2 h-4 w-4" />
                  Roadmap
                </Link>
                <Link
                  href="/dashboard/pricing"
                  className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                >
                  <CreditCard className="mr-2 h-4 w-4" />
                  Pricing
                </Link>
                {user?.role === "ADMIN" && (
                  <>
                    <Link
                      href="/dashboard/admin/quality-report"
                      className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                    >
                      <BarChart3 className="mr-2 h-4 w-4" />
                      Qualidade
                    </Link>
                    <Link
                      href="/dashboard/admin/settings"
                      className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                    >
                      <ShieldCheck className="mr-2 h-4 w-4" />
                      Admin
                    </Link>
                    <Link
                      href="/dashboard/admin/roadmap"
                      className="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
                    >
                      <Map className="mr-2 h-4 w-4" />
                      Roadmap Admin
                    </Link>
                  </>
                )}
              </nav>
            </div>

            <div className="flex items-center gap-2">
              <NotificationBell />
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                    <Avatar className="h-8 w-8">
                      <AvatarFallback className="bg-primary/20 text-primary">
                        {user?.name?.charAt(0).toUpperCase() || "U"}
                      </AvatarFallback>
                    </Avatar>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-56" align="end">
                  <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col space-y-1">
                      <p className="text-sm font-medium leading-none">{user?.name}</p>
                      <p className="text-xs leading-none text-muted-foreground">
                        {user?.email}
                      </p>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem>
                    <Link href="/dashboard/profile" className="cursor-pointer flex items-center w-full">
                      <UserIcon className="mr-2 h-4 w-4" />
                      <span>Meu Perfil</span>
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    className="cursor-pointer text-destructive focus:bg-destructive/10 focus:text-destructive"
                    onClick={() => logout()}
                  >
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Sair</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </header>

        {/* Main Content */}
        <main className="flex-1 container mx-auto p-4 md:p-8 animate-in fade-in duration-500">
          {children}
        </main>
      </div>
    </PrivateRoute>
  );
}
