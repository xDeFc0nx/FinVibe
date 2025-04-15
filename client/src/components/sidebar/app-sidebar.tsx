import * as React from "react";
import { LayoutDashboard, CreditCard } from "lucide-react";

import { NavUser } from "@/components/sidebar/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { NavMain } from "@/components/sidebar/nav-main";
import { AccountSwitcher } from "@/components/sidebar/account-switcher";
import { useSelector } from 'react-redux';
import type { RootState } from '@/store/store.ts';
const data = {
  Links: [
    {
      name: "Dashboard",
      url: "/app/dashboard",
      icon: LayoutDashboard,
    },
    {
      name: "Transactions",
      url: "/app/transactions",
      icon: CreditCard,
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {

  const userData = useSelector((state: RootState) => state.user.data)
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <AccountSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.Links} />
      </SidebarContent>
      <SidebarFooter>
         {userData && <NavUser user={userData} />}   
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
