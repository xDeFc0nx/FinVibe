import { useUserData } from '@/components/context/userData';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { SidebarInset, SidebarTrigger } from '@/components/ui/sidebar';
import { DataTable } from './dataTable';
import { columns } from './columns';
import { BalanceChart, ExpensesChart, IncomeChart } from '@/components/charts/Charts';
import { useWebSocket } from '@/components/WebSocketProvidor';
import { useState } from 'react';

export default function Index() {
  const {
    setTransactions,
    transactions,
    activeAccount,
    setAccounts,
    setActiveAccount,
    dateRange,
    setDateRange,
    chartOverview,
    setChartOverview,
  } = useUserData();

  const { socket } = useWebSocket();
  const handleDateRangeChange = (value: string) => {
       setDateRange(value);

      if (socket && activeAccount) {
      const updatedAccount = { ...activeAccount, DateRange: value };
      setActiveAccount(updatedAccount);

           const messageHandler = (msg: string) => {
        const response = JSON.parse(msg);

        if (response.transactions) {
          setTransactions(response.transactions);
        }

        if (response.totalIncome !== undefined) {
          setAccounts((prev) =>
            prev.map((acc) =>
              acc.AccountID === activeAccount.AccountID
                ? { ...acc, Income: response.totalIncome }
                : acc,
            ),
          );
          setActiveAccount((prev) =>
            prev && prev.AccountID === activeAccount.AccountID
              ? { ...prev, Income: response.totalIncome }
              : prev,
          );
        }

        if (response.totalExpense !== undefined) {
          setAccounts((prev) =>
            prev.map((acc) =>
              acc.AccountID === activeAccount.AccountID
                ? { ...acc, Expense: response.totalExpense }
                : acc,
            ),
          );
          setActiveAccount((prev) =>
            prev && prev.AccountID === activeAccount.AccountID
              ? { ...prev, Expense: response.totalExpense }
              : prev,
          );
        }

        if (response.accountBalance !== undefined) {
          setAccounts((prev) =>
            prev.map((acc) =>
              acc.AccountID === activeAccount.AccountID
                ? { ...acc, AccountBalance: response.accountBalance }
                : acc,
            ),
          );
          setActiveAccount((prev) =>
            prev && prev.AccountID === activeAccount.AccountID
              ? { ...prev, AccountBalance: response.accountBalance }
              : prev,
          );
        }
        if(response.chartData){
            setChartOverview(response.chartData)
        }

      };

      socket.onMessage(messageHandler);
    }
  };

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
          <Select
            value={activeAccount?.DateRange || 'this_month'}
            onValueChange={handleDateRangeChange}
          >
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
        <div className="grid auto-rows-min gap-4 md:grid-cols-4">
          <div className="rounded-xl border bg-card text-card-foreground shadow">
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">Balance</h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">
                {activeAccount?.AccountBalance}
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
              <div className="text-2xl font-bold"> {activeAccount?.Income}</div>
              <IncomeChart/>
                          </div>
          </div>
          <div className="rounded-xl border bg-card text-card-foreground shadow">
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">
                Total Expenses
              </h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">{activeAccount?.Expense}</div>
              <ExpensesChart/>
                        </div>
          </div>
          <div className="aspect-video rounded-xl bg-muted/50" />
        </div>
        <div className="grid auto-rows-min gap-4 md:grid-cols-2">
          <div className="rounded-xl bg-muted/50" />
          <DataTable columns={columns} data={transactions} />
        </div>
      </div>
    </SidebarInset>
  );
}
