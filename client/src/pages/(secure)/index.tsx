import { useUserData } from '@/components/context/userData';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';

import { Separator } from '@/components/ui/separator';
import { SidebarInset, SidebarTrigger } from '@/components/ui/sidebar';
import { DataTable } from './dataTable';
import { columns } from './columns';
import { BalanceChart } from '@/components/charts/balanceChart';
export default function Index() {
  const { activeAccount, transactions } = useUserData();
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
      </header>
      <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
        <div className="grid auto-rows-min gap-4 md:grid-cols-4">
          <div className="rounded-xl border bg-card text-card-foreground shadow">
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">Balance</h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">
                {' '}
                {activeAccount?.AccountBalance}
              </div>
          <BalanceChart/> 
            </div>
          </div>
          <div className="rounded-xl border bg-card text-card-foreground shadow">
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">Total Income</h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">
                {' '}
                {activeAccount?.Income}
              </div>
              <p className="text-xs text-muted-foreground">
                +20.1% from last month
              </p>
            </div>
          </div>
         <div className="rounded-xl border bg-card text-card-foreground shadow">
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">Total Expenses</h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">
                {' '}
                {activeAccount?.Expense}
              </div>
              <p className="text-xs text-muted-foreground">
                +20.1% from last month
              </p>
            </div>
          </div>
          <div className="aspect-video rounded-xl bg-muted/50" />
        </div>
        <div className="grid auto-rows-min gap-4 md:grid-cols-2">
          <div className=" rounded-xl bg-muted/50" />

          <DataTable columns={columns} data={transactions} />
        </div>
      </div>
    </SidebarInset>
  );
}
