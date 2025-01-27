import * as z from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { toast } from 'react-toastify';
import { useUserData } from '@/components/context/userData';
import { useWebSocket } from '@/components/WebSocketProvidor';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { AccountSwitcher } from '../sidebar/account-switcher';

const formSchema = z.object({
  Type: z.enum(['Income', 'Expense']),
  Description: z.string().min(1, 'Description is required'),
  Amount: z.number().min(1, 'Amount must be greater than 0'),
  IsRecurring: z.boolean(),
});

export const AddTransaction = () => {
  const { socket, isReady } = useWebSocket();
  const {
    setTransactions,
    activeAccount,
    setAccounts,
    setActiveAccount,
    setChartOverview,
    dateRange,
  } = useUserData();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      Type: 'Income',
      Description: '',
      Amount: 0,
      IsRecurring: false,
    },
  });

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    try {
      if (socket && isReady && activeAccount) {
        const currentAccountId = activeAccount.AccountID;

        socket.send('createTransaction', {
          AccountID: currentAccountId,
          ...values,
        });

        const balanceHandler = (msg: string) => {
          const response = JSON.parse(msg);

          if (response.transaction) {
            setTransactions((prev) => [...prev, response.transaction]);
            toast.success('Transaction added!', {
              toastId: 'success',
            });
            form.reset();
          }

          if (response.totalIncome !== undefined) {
            setAccounts((prev) =>
              prev.map((acc) =>
                acc.AccountID === currentAccountId
                  ? { ...acc, Income: response.totalIncome }
                  : acc,
              ),
            );
          }

          if (response.totalExpense !== undefined) {
            setAccounts((prev) =>
              prev.map((acc) =>
                acc.AccountID === currentAccountId
                  ? { ...acc, Expense: response.totalExpense }
                  : acc,
              ),
            );
          }

          if (response.accountBalance !== undefined) {
            setAccounts((prev) =>
              prev.map((acc) =>
                acc.AccountID === currentAccountId
                  ? { ...acc, AccountBalance: response.accountBalance }
                  : acc,
              ),
            );
            setActiveAccount((prev) =>
              prev?.AccountID === currentAccountId
                ? { ...prev, AccountBalance: response.accountBalance }
                : prev,
            );
          }
          if (response.chartData) {
            setChartOverview(response.chartData);
          }
        };

        socket.onMessage(balanceHandler);
        setTimeout(() => {
          socket.send('getAccountIncome', { AccountID: currentAccountId });
          socket.send('getAccountExpense', { AccountID: currentAccountId });
          socket.send('getAccountBalance', { AccountID: currentAccountId });
          socket.send('getCharts', {
            AccountID: currentAccountId,
            DataRange: dateRange,
          });
        }, 100);
        return;
      }
    } catch (error) {
      console.error('Submission error', error);
      toast.error('Failed to add transaction');
    }
  };

  return (
    <Dialog>
      <DialogTrigger>
        <Button variant="green" className="pl-5">
          New
        </Button>
      </DialogTrigger>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add a New Transaction</DialogTitle>
          <DialogDescription>
            Fill out the details to add a new transaction.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="Type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Transaction Type</FormLabel>
                  <FormControl>
                    <Select
                      onValueChange={field.onChange}
                      value={field.value}
                      defaultValue={field.value}
                    >
                      <SelectTrigger className="rounded-lg h-12">
                        <SelectValue placeholder="Select transaction type" />
                      </SelectTrigger>
                      <SelectContent className="rounded-lg min-w-[200px]">
                        <SelectItem value="Income" className="cursor-pointer">
                          Income
                        </SelectItem>
                        <SelectItem value="Expense" className="cursor-pointer">
                          Expense
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="Description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Input placeholder="Description" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="Amount"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Amount</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Amount"
                      type="number"
                      {...field}
                      onChange={(e) => field.onChange(Number(e.target.value))}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="IsRecurring"
              render={({ field }) => (
                <FormItem className="flex items-center space-x-2">
                  <FormControl>
                    <Checkbox
                      checked={field.value}
                      onCheckedChange={(checked) => field.onChange(!!checked)}
                      ref={field.ref}
                    />
                  </FormControl>
                  <FormLabel>Recurring</FormLabel>
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button variant="default" type="submit">
                Add!
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};
