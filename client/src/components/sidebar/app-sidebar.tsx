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
import { useUserData } from "@/components/context/userData";
import { NavMain } from "@/components/sidebar/nav-main";
import { AccountSwitcher } from "@/components/sidebar/account-switcher";

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
	const { userData } = useUserData();

	return (
		<Sidebar collapsible="icon" {...props}>
			<SidebarHeader>
				<AccountSwitcher />
			</SidebarHeader>
			<SidebarContent>
				<NavMain items={data.Links} />
			</SidebarContent>
			<SidebarFooter>
				<NavUser user={userData} />
			</SidebarFooter>
			<SidebarRail />
		</Sidebar>
	);
}
