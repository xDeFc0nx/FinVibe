import { WebSocketProvider } from "@/components/WebSocketProvidor";
import CheckAuth from "@/lib/checkAuth";
import type { JSX } from "react";
import { Outlet } from "react-router-dom";

interface LayoutProps {
  children?: React.ReactNode;
}

export function Layout({ children }: LayoutProps): JSX.Element {
  return (
    <main>
      <CheckAuth />
      <WebSocketProvider />
      <Outlet />
    </main>
  );
}
