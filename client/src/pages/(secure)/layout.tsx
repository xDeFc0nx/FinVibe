import { WebSocketProvider } from "@/components/WebSocketProvidor";
import { AppSidebar } from "@/components/app-sidebar";
import { UserDataProvider } from "@/components/context/userData";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import CheckAuth from "@/lib/checkAuth";
import type { JSX } from "react";
import type React from "react";
import { Outlet } from "react-router-dom";
interface LayoutProps {
  children?: React.ReactNode;
}

export function Layout({ children }: LayoutProps): JSX.Element {
  return (
    <main>
      <CheckAuth />
      <WebSocketProvider>
       <UserDataProvider>

        <SidebarProvider>
          <AppSidebar />
          <Outlet />
        </SidebarProvider>
       </UserDataProvider>
      </WebSocketProvider>
    </main>
  );
}
