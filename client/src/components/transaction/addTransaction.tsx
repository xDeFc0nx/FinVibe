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

const formSchema = z.object({
  Description: z.string().min(1, 'Description is required'),
  Amount: z.number().min(1, 'Amount must be greater than 0'),
  IsRecurring: z.boolean(),
});

export const AddTransaction = () => {
  const { socket, isReady } = useWebSocket();
  const { transactions, setTransactions, activeAccount } = useUserData();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      Description: '',
      Amount: 0,
      IsRecurring: false,
    },
  });

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    try {
      if (socket && isReady) {
        socket.send('createTransaction', {
          AccountID: activeAccount?.AccountID,
          Description: values.Description,
          Amount: values.Amount,
          IsRecurring: values.IsRecurring,
        });

        socket.onMessage((msg) => {
          const response = JSON.parse(msg);

          if (response.transaction) {
            setTransactions([
              ...transactions,
              {
                ID: response.transaction.ID,
                UserID: response.transaction.UserID,
                AccountID: response.transaction.AccountID,
                Description: values.Description,
                Amount: values.Amount,
                IsRecurring: values.IsRecurring,
                CreatedAt: response.transaction.CreatedAt,
              },
            ]);
            toast.success('Transaction added successfully!');
          }
          if (response.Error) {
            toast.error(response.Error);
          }
        });
      }
    } catch (error) {
      console.error('Form submission error', error);
      toast.error('Failed to submit the form. Please try again.');
    }
  };

  return (
    <Dialog>
      <DialogTrigger>
        <Button variant="default" className="pl-5">
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
