import * as React from "react";
import { ChevronsUpDown, Plus } from "lucide-react";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuShortcut,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	useSidebar,
} from "@/components/ui/sidebar";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { toast } from "react-toastify";
import { useUserData } from "@/components/context/userData";
import { useWebSocket } from "@/components/WebSocketProvidor";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { isRegExp } from "util/types";

export function AccountSwitcher() {
	const {
		accounts,
		setAccounts,
		activeAccount,
		setActiveAccount,
		setTransactions,
		dateRange,
		setChartOverview,
		setIncomePie,
		setExpensesPie,
		refresh,
	} = useUserData();
	const { isMobile } = useSidebar();
	const [open, setOpen] = React.useState(false);

	const { socket, isReady } = useWebSocket();
  console.log("AccountSwitcher RENDER: activeAccount is:", JSON.stringify(activeAccount));	const formSchema = z.object({
		Type: z.string().min(1, "Account type is required"),
	});

	const form = useForm({
		resolver: zodResolver(formSchema),
		defaultValues: {
			Type: "",
		},
	});
	function saveAccount(account: any) {
		setActiveAccount(account);
		localStorage.setItem("activeAccount", JSON.stringify(account));
		console.log("Saved Account to LocalStorage:", account);
	}
	function onSubmit(values: z.infer<typeof formSchema>) {
		try {
			console.log(values);

			if (socket && isReady) {
				socket.send("createAccount", {
					Type: values.Type,
				});

				socket.onMessage((msg) => {
					const response = JSON.parse(msg);

					if (response.account) {
						setOpen(false);
						saveAccount(response.account);
						setAccounts((prevAccounts) => [...prevAccounts, response.account]);
						setActiveAccount(response.account);
					}

					if (response.Error) {
						toast.error(response.Error);
					}
				});
			}
		} catch (error) {
			console.error("Form submission error", error);
			toast.error("Failed to submit the form. Please try again.");
		}
	}

	React.useEffect(() => {
		if (socket && isReady && activeAccount?.id) {
			const currentAccountId = activeAccount.id;

			socket.send("getTransactions", {
				AccountID: activeAccount?.id,
				DateRange: dateRange,
			});
			socket.send("getAccountIncome", {
				AccountID: activeAccount?.id,
				DateRange: dateRange,
			});
			socket.send("getAccountExpense", {
				AccountID: activeAccount?.id,
				DateRange: dateRange,
			});
			socket.send("getAccountBalance", {
				AccountID: activeAccount?.id,
				DateRange: dateRange,
			});
			socket.send("getCharts", {
				AccountAccountID: activeAccount?.id,
				DateRange: dateRange,
			});

			socket.onMessage((msg) => {
				const response = JSON.parse(msg);
				if (response.transactions) {
					setTransactions(response.transactions);
				}

				if (response.Error) {
					toast.error(response.Error);
				}
				if (response.totalIncome !== undefined) {
					setAccounts((prev) =>
						prev.map((acc) =>
							acc.id === currentAccountId
								? { ...acc, Income: response.totalIncome }
								: acc,
						),
					);
				}

				if (response.totalExpense !== undefined) {
					setAccounts((prev) =>
						prev.map((acc) =>
							acc.id === currentAccountId
								? { ...acc, Expense: response.totalExpense }
								: acc,
						),
					);
				}

				if (response.accountBalance !== undefined) {
					setAccounts((prev) =>
						prev.map((acc) =>
							acc.id === currentAccountId
								? { ...acc, AccountBalance: response.accountBalance }
								: acc,
						),
					);
					setActiveAccount((prev) =>
						prev?.id === currentAccountId
							? { ...prev, AccountBalance: response.accountBalance }
							: prev,
					);
				}
				if (response.chartData) {
					setChartOverview(response.chartData);
				}
				if (response.IncomePie) {
					setIncomePie(response.IncomePie);
				}
				if (response.ExpensesPie) {
					setExpensesPie(response.ExpensesPie);
				}
			});
		}
	}, [activeAccount?.id, dateRange, isReady, refresh]);
	return (
		<SidebarMenu>
			<Dialog open={open} onOpenChange={setOpen}>
				<SidebarMenuItem>
					<DropdownMenu>
						<DropdownMenuTrigger asChild>
							<SidebarMenuButton
								size="lg"
								className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
							>
								<div className="grid flex-1 text-left text-sm leading-tight">
									<span className="truncate font-semibold">
										{activeAccount?.type}
									</span>
								</div>
								<ChevronsUpDown className="ml-auto" />
							</SidebarMenuButton>
						</DropdownMenuTrigger>
						<DropdownMenuContent
							className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
							align="start"
							side={isMobile ? "bottom" : "right"}
							sideOffset={4}
						>
							<DropdownMenuLabel className="text-xs text-muted-foreground">
								Accounts
							</DropdownMenuLabel>
							{accounts.map((account) => (
								<DropdownMenuItem
									key={account.id}
									onClick={() => saveAccount(account)}
									className="gap-2 p-2"
								>
									{account.type}
									<DropdownMenuShortcut>
										⌘{accounts.indexOf(account) + 1}
									</DropdownMenuShortcut>
								</DropdownMenuItem>
							))}
							<DropdownMenuSeparator />
							<DropdownMenuItem className="gap-2 p-2">
								<div className="flex size-6 items-center justify-center rounded-md border bg-background">
									<Plus className="size-4" />
								</div>
								<DialogTrigger>Add account</DialogTrigger>
							</DropdownMenuItem>
						</DropdownMenuContent>
					</DropdownMenu>
				</SidebarMenuItem>

				<DialogContent>
					<DialogHeader>
						<DialogTitle>Add New Account</DialogTitle>
						<DialogDescription>
							<Form {...form}>
								<form
									onSubmit={form.handleSubmit(onSubmit)}
									className="space-y-4"
								>
									<FormField
										control={form.control}
										name="Type"
										render={({ field }) => (
											<FormItem>
												<FormLabel>Account Type</FormLabel>
												<FormControl>
													<Input
														placeholder="Account Type"
														type="text"
														{...field}
													/>
												</FormControl>
												<FormMessage />
											</FormItem>
										)}
									/>
									<DialogFooter>
										<Button variant="default" className="mt-5" type="submit">
											Add!
										</Button>
									</DialogFooter>
								</form>
							</Form>
						</DialogDescription>
					</DialogHeader>
				</DialogContent>
			</Dialog>
		</SidebarMenu>
	);
}
