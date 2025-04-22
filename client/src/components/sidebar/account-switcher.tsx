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
import { useSelector, useDispatch } from 'react-redux';
import type { RootState, AppDispatch } from '@/store/store.ts';
import { accountsReceived, setActiveAccount, updateAccountDetails } from '@/store/slices/accountsSlice';
import { transactionsReceived } from '@/store/slices/transactionsSlice';
import { overviewReceived } from '@/store/slices/chartsSlice'
import type { Account, } from '@/types';

export function AccountSwitcher() {

  const { isMobile } = useSidebar();
  const [open, setOpen] = React.useState(false);

  const dispatch: AppDispatch = useDispatch();
  const { socket, isReady } = useWebSocket();

  const formSchema = z.object({
    Type: z.string().min(1, "Account type is required"),
  });

  const dateRange = useSelector((state: RootState) => state.transactions.dateRange);
  const { activeAccountId, list: currentAccounts } = useSelector((state: RootState) => state.accounts);
  const activeAccount: Account | null = React.useMemo(() => {
    if (!activeAccountId) return null;
    return currentAccounts.find(acc => acc.id === activeAccountId) || null;
  }, [activeAccountId, currentAccounts]);
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      Type: "",
    },
  });
  function saveAccount(account: any) {
    dispatch(setActiveAccount(account.id));
    localStorage.setItem("activeAccount", JSON.stringify(account));
  }
  function onSubmit(values: z.infer<typeof formSchema>) {
    try {

      if (socket && isReady) {
        socket.send("createAccount", {
          type: values.Type,
        });

        socket.onMessage((msg) => {
          const response = JSON.parse(msg);

          if (response.account) {
            const newAccount: Account = response.account;
            setOpen(false);
            saveAccount(newAccount);
            const updatedAccountList: Account[] = [...currentAccounts, newAccount];
            dispatch(accountsReceived(updatedAccountList));

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
    if (socket && isReady && activeAccountId) {

      socket.send("getTransactions", {
        AccountID: activeAccountId,
        DateRange: dateRange,
      });
      socket.send("getAccount", {
        AccountID: activeAccountId,
        // DateRange: dateRange,
      });
      socket.send("getCharts", {
        AccountID: activeAccountId,
        DateRange: dateRange,
      });

      socket.onMessage((msg) => {

        const response = JSON.parse(msg);
        if (response.transactions) {
          dispatch(transactionsReceived(response.transactions))
        }

        if (response.Error) {
          toast.error(response.Error);
        }
        if (response.AccountData) {
          dispatch(updateAccountDetails({
            id: activeAccountId,
            details: response.AccountData,
          }));
        }
        if (response.chartData && response.IncomePie && response.ExpensesPie) {
          const overviewPayload = {
            chart: response.chartData,
            incomePie: response.IncomePie,
            expensePie: response.ExpensesPie
          };

          dispatch(overviewReceived(overviewPayload));

        }

      });
    }
  }, [activeAccountId, dateRange, isReady,]);
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
              {currentAccounts.map((account) => (
                <DropdownMenuItem
                  key={activeAccountId}
                  onClick={() => saveAccount(account)}
                  className="gap-2 p-2"
                >
                  {account.type}
                  <DropdownMenuShortcut>
                    ⌘{currentAccounts.indexOf(account) + 1}
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
