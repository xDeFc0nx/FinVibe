import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { SidebarInset, SidebarTrigger } from "@/components/ui/sidebar";
import { DataTable } from "./dataTable";
import { columns } from "./columns";
import {
  BalanceChart,
  ExpensesChart,
  ExpensesPie,
  IncomeChart,
  IncomePie,
} from "@/components/charts/Charts";
import { useWebSocket } from "@/components/WebSocketProvidor";
import { useEffect } from "react";
import { useSelector, useDispatch } from 'react-redux';
import type { RootState, AppDispatch } from '@/store/store.ts';
import { accountsReceived, updateAccountDetails, setActiveAccount } from '@/store/slices/accountsSlice';
import { setDateRange, transactionsReceived, } from '@/store/slices/transactionsSlice';
import { overviewReceived, } from '@/store/slices/chartsSlice';
import { userReceived } from "@/store/slices/userSlice";

export default function Index() {
  const dispatch: AppDispatch = useDispatch();
  const { list: transactions, dateRange } = useSelector((state: RootState) => state.transactions);
  const { list: accounts, activeAccountId, } = useSelector((state: RootState) => state.accounts);
  const userData = useSelector((state: RootState) => state.user.data);
  const activeAccount = accounts.find((acc) => acc.id === activeAccountId) || null;
  const { socket, isReady } = useWebSocket();
  useEffect(() => {
    const storedAccountJSON = localStorage.getItem("activeAccount");
    if (socket && isReady) {


      socket.send("getUser");
      socket.send("getAccounts");
      socket.onMessage((msg) => {
        const response = JSON.parse(msg);
        if (response.accounts) {

          dispatch(accountsReceived(response.accounts))
        }
        if (response.userData) {
          dispatch(userReceived(response.userData))
        }
        if (response.AccountData){
          dispatch(updateAccountDetails(response.AccountData))
        }

      });
    }

    if (storedAccountJSON) {
      try {
        const storedAccount = JSON.parse(storedAccountJSON);
        const accountIdFromStorage = storedAccount.ID ?? storedAccount.id;
        if (typeof accountIdFromStorage === 'string' && accountIdFromStorage.length > 0) {
          dispatch(setActiveAccount(accountIdFromStorage));
        } else {
          localStorage.removeItem("activeAccount");
        }
      } catch (error) {
        console.error("   Error parsing stored activeAccount from localStorage:", error);
        localStorage.removeItem("activeAccount");
      }
    }


  }, [socket, isReady, dispatch]);

  const handleDateRangeChange = (value: string) => {
    setDateRange(value);
    const activeAccount = accounts.find(acc => acc.id === activeAccountId) || null;
    if (socket && activeAccount) {

      const messageHandler = (msg: string) => {
        const response = JSON.parse(msg);

        if (response.transactions) {
          dispatch(transactionsReceived(response.transactions));
        }
        if (response.AccountData) {
          if (activeAccountId) {

            dispatch(updateAccountDetails({
              id: activeAccountId,
              details: response.AccountData,
            }));
          }
        }
        if (response.chartData) {
          dispatch(overviewReceived(response.chartData));
        }

        socket.onMessage(messageHandler);
      }
    };
  }

  return (
    <SidebarInset>
      <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
        <div className="flex items-center gap-2 px-4">
          <SidebarTrigger className="-ml-1" />
          <Separator orientation="vertical" className="mr-2 h-4" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbSeparator className="hidden md:block" />
              <BreadcrumbItem>
                <BreadcrumbPage>Dashboard</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </div>
        <div className="flex-1 flex justify-start">
          <Select value={dateRange} onValueChange={handleDateRangeChange}>
            {" "}
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Select Date Range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="this_month">This Month</SelectItem>
              <SelectItem value="last_month">Last Month</SelectItem>
              <SelectItem value="last_6_months">Last 6 Months</SelectItem>
              <SelectItem value="this_year">This Year</SelectItem>
              <SelectItem value="last_year">Last Year</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </header>
      <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
        <div className="grid auto-rows-min gap-4 md:grid-cols-3">
          <div className="rounded-xl border bg-card text-card-foreground shadow">
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">Balance</h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">
                {userData?.Currency}
                {activeAccount?.balance}
              </div>

              <BalanceChart />
            </div>
          </div>
          <div className="rounded-xl border bg-card text-card-foreground shadow">
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">
                Total Income
              </h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">
                {" "}
                {userData?.Currency}
                {activeAccount?.income}
              </div>
              <IncomeChart />
            </div>
          </div>
          <div className="rounded-xl border bg-card text-card-foreground shadow">
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">
                Total Expenses
              </h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">
                {" "}
                {userData?.Currency}
                {activeAccount?.expense}
              </div>
              <ExpensesChart />
            </div>
          </div>
        </div>
        <div className="grid auto-rows-min gap-4 md:grid-cols-3">
          <IncomePie />
          <ExpensesPie />
          <DataTable columns={columns} data={transactions} />
        </div>
      </div>
    </SidebarInset>
  );

} 
