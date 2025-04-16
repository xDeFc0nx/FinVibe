import * as React from "react";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Descriptions } from "@/components/Descriptions";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { cn } from "@/lib/utils";
import { CalendarIcon } from "lucide-react";
import { Calendar } from "../ui/calendar";
import { format } from "date-fns";
import { useSelector, useDispatch } from 'react-redux';
import type { RootState, AppDispatch } from '@/store/store.ts';
import { addTransaction } from '@/store/slices/transactionsSlice';
import { updateAccountDetails} from "@/store/slices/accountsSlice"
import type { Account } from '@/types';
const formSchema = z.object({
  Type: z.enum(["Income", "Expense"]),
  Description: z.string().min(1, "Description is required"),
  Amount: z.number().min(1, "Amount must be greater than 0"),
  IsRecurring: z.boolean(),
  CreatedAt: z.date(),
});

export const AddTransaction = () => {
  const dispatch: AppDispatch = useDispatch();
  const { socket, isReady } = useWebSocket();
  const activeAccountId = useSelector((state: RootState) => state.accounts.activeAccountId);
  const currentAccounts = useSelector((state: RootState) => state.accounts.list);
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      Type: "Income",
      Description: "",
      CreatedAt: new Date(),
      IsRecurring: false,
    },
  });
  const activeAccount: Account | null = React.useMemo(() => {
    if (!activeAccountId) return null;
    return currentAccounts.find(acc => acc.id === activeAccountId) || null;
  }, [activeAccountId, currentAccounts]);


  const onSubmit = (values: z.infer<typeof formSchema>) => {
    try {
      if (socket && isReady && activeAccount) {
        const payload = {
          ...values,
          CreatedAt: values.CreatedAt.toISOString(),
        };
        socket.send("createTransaction", {
          AccountID: activeAccountId,
          ...payload,
        });

        socket.onMessage((msg) => {
          const response = JSON.parse(msg);
          if (response.transaction) {
            dispatch(addTransaction(response.Transaction))
            toast.success("Transaction added!", { toastId: "success" });
            form.reset();
          }

          if (response.AccountData) {
            if (activeAccountId) {


              dispatch(updateAccountDetails({
                id: activeAccountId,
                details: response.AccountData,
              }));
            }else{
              console.log("account id not found")
            }

          }
        });
      }
    } catch (error) {
      console.error("Submission error", error);
      toast.error("Failed to add transaction");
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
              render={({ field }) => {
                const transactionType = form.watch("Type");
                return (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                        defaultValue={field.value}
                      >
                        <SelectTrigger className="rounded-lg h-12">
                          <SelectValue placeholder="Select description" />
                        </SelectTrigger>
                        <SelectContent className="rounded-lg min-w-[200px]">
                          {Descriptions[transactionType].map((desc) => (
                            <SelectItem
                              key={desc}
                              value={desc}
                              className="cursor-pointer"
                            >
                              {desc}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                );
              }}
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
              name="CreatedAt"
              render={({ field }) => (
                <FormItem className="flex flex-col">
                  <FormLabel>Date</FormLabel>
                  <Popover>
                    <PopoverTrigger asChild>
                      <FormControl>
                        <Button
                          variant={"outline"}
                          className={cn(
                            "w-[240px] pl-3 text-left font-normal",
                            !field.value && "text-muted-foreground",
                          )}
                        >
                          {field.value ? (
                            format(field.value, "PPP")
                          ) : (
                            <span>Pick a date</span>
                          )}
                          <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                        </Button>
                      </FormControl>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                      <Calendar
                        mode="single"
                        selected={field?.value}
                        onSelect={field.onChange}
                        disabled={(date) =>
                          date > new Date() || date < new Date("1900-01-01")
                        }
                      />
                    </PopoverContent>
                  </Popover>
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
