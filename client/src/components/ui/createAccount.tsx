"use client";
import * as React from "react";
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
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import {  useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import * as z from "zod";
import { useWebSocket } from "@/components/WebSocketProvidor";
import { useEffect } from "react";
import { useSelector, useDispatch } from 'react-redux';
import type { RootState, AppDispatch } from '@/store/store.ts';
import { addAccount, setActiveAccount } from '@/store/slices/accountsSlice';
import type { Account } from '@/types'; // Adjust path
const formSchema = z.object({
  Type: z.string().min(1, "Account type is required"),
});

export default function CreateAccount() {
  // const navigate = useNavigate();
  const dispatch: AppDispatch = useDispatch();
  const { activeAccountId, list: currentAccounts } = useSelector((state: RootState) => state.accounts);
  const { socket, isReady } = useWebSocket();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });
  const activeAccount = React.useMemo(() => {
    if (!activeAccountId) return null;
    return currentAccounts.find(acc => acc.ID === activeAccountId) || null;
  }, [activeAccountId, currentAccounts]);
  function handleSubmit(values: z.infer<typeof formSchema>) {
    try {
      console.log(values);

      if (socket && isReady) {
        socket.send("createAccount", {
          Type: values.Type,
        });

        socket.onMessage((msg) => {
          const response = JSON.parse(msg);

          if (response.account) {
            console.log(response.account);
            const newAccount: Account = response.account;
            dispatch(addAccount(newAccount));
            dispatch(setActiveAccount(newAccount.ID));
            // navigate("/app/dashboard");
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

  useEffect(() => {
    if (activeAccount) {
      localStorage.setItem("activeAccount", JSON.stringify(activeAccount));
    }
  }, [activeAccount]);
  return (
    <>
      <div className=" 0">
        <div className="lg:p-8">
          <div className="flex flex-col space-y-2 text-center">
            <h1 className="text-2xl font-semibold tracking-tight">
              Create Account
            </h1>
            <p className="text-sm text-muted-foreground">
              This is where all your transactions are stored,
            </p>
          </div>

          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              <FormField
                control={form.control}
                name="Type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Account Type</FormLabel>
                    <FormControl>
                      <Input placeholder="Account Type" type="" {...field} />
                    </FormControl>

                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" className="mt-5">
                Create
              </Button>
            </form>
          </Form>
        </div>
      </div>
    </>
  );
}
