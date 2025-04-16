import { AppSidebar } from "@/components/sidebar/app-sidebar";
import { SidebarProvider } from "@/components/ui/sidebar";
import { WebSocketProvider } from "@/components/WebSocketProvidor";
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
        <SidebarProvider>
          <AppSidebar />
          <Outlet />
        </SidebarProvider>
      </WebSocketProvider>

    </main>
  );
}
