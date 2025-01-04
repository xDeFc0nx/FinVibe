"use client";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import * as z from "zod";

const formSchema = z.object({
  email: z.string().email("Invalid email format"),
  password: z.string().min(6, "Password must be at least 8 characters long"),
});

export default function MyForm() {
  const navigate = useNavigate();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  const handleSubmit = async (data: z.infer<typeof formSchema>) => {
    try {
      const response = await fetch("http://localhost:3001/Login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
        credentials: "include",
      });
      const responseData = await response.json();
      console.log(responseData);

      if (response.ok) {
        document.cookie = `jwt=${responseData.token}; path=/; secure; httpOnly`;
        navigate("/app/dashboard");
      } else {
        toast("Login Failed Try again");
      }
    } catch (error) {
      toast("Login Failed Try again");
    }
  };

  return (
    <>
      <div className="container relative hidden h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-1 lg:px-0">
        <Link to="/register">
          <Button className="absolute right-4 top-4 md:right-8 md:top-8">
            Register
          </Button>
        </Link>

        <div className="lg:p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <div className="flex flex-col space-y-2 text-center">
              <h1 className="text-2xl font-semibold tracking-tight">
                {" "}
                Welcome Back!
              </h1>
              <p className="text-sm text-muted-foreground">
                Your data has been carefully wrapped in a secure layer of
                encryption and stored in a digital vault,
              </p>
            </div>

            <Form {...form}>
              <form onSubmit={form.handleSubmit(handleSubmit)}>
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input
                          className="grid w-full max-w-sm items-center gap-1.5"
                          placeholder="Email"
                          type="email"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>Enter your email</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Password</FormLabel>
                      <FormControl>
                        <Input
                          className="grid w-full max-w-sm items-center gap-1.5"
                          placeholder="Password"
                          type="password"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>Enter your password</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button type="submit">Submit</Button>
              </form>
            </Form>
            <p className="px-8 text-center text-sm text-muted-foreground pt-5">
              By clicking continue, you agree to our{" "}
              <a
                href="/terms"
                className="underline underline-offset-4 hover:text-primary"
              >
                Terms of Service
              </a>{" "}
              and{" "}
              <a
                href="/privacy"
                className="underline underline-offset-4 hover:text-primary"
              >
                Privacy Policy
              </a>
              .
            </p>
          </div>
        </div>
      </div>
    </>
  );
}
