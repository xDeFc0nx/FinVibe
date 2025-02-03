"use client";
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
import { PasswordInput } from "@/components/ui/password-input";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { toast } from "react-toastify";
import * as z from "zod";

const formSchema = z.object({
	email: z.string().email("Invalid email format"),
	password: z.string().min(8, "Password must be at least 8 characters long"),
});
export default function Login() {
	const navigate = useNavigate();
	const [wrongCredentials, setWrongCredentials] = useState(false);

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
	});
	const onSubmit = async (data: z.infer<typeof formSchema>) => {
		const success = await handleLogin(data, setWrongCredentials);
		if (success) navigate("/app/dashboard");
	};
	return (
		<>
			<div className="lg:p-8">
				<div className="flex flex-col space-y-2 text-center">
					<h1 className="text-2xl font-semibold tracking-tight">
						Welcome Back!
					</h1>
					<p className="text-sm text-muted-foreground">
						Your data has been carefully wrapped in a secure layer of encryption
						and stored in a digital vault,
					</p>
				</div>
				<br />
				<Form {...form}>
					<form onSubmit={form.handleSubmit(onSubmit)}>
						<FormField
							control={form.control}
							name="email"
							render={({ field }) => (
								<FormItem>
									<FormControl>
										<Input
											className="grid w-full max-w-sm items-center gap-1.5"
											placeholder="Email"
											type="email"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<br />
						<FormField
							control={form.control}
							name="password"
							render={({ field }) => (
								<FormItem>
									<FormControl>
										<PasswordInput
											className="grid w-full max-w-sm items-center gap-1.5"
											placeholder="Password"
											type="password"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						{wrongCredentials ? (
							<p className="px-8 text-center text-sm text-red-500  pt-5">
								Wrong Email or Password
							</p>
						) : (
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
						)}

						<Button type="submit" className="mt-5">
							Login
						</Button>
					</form>
				</Form>
			</div>{" "}
		</>
	);
}
export const handleLogin = async (
	data: {
		email: string;
		password: string;
	},
	setWrongCredentials: (value: boolean) => void,
) => {
	try {
		const response = await fetch("http://localhost:3001/Login", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(data),
			credentials: "include",
		});
		const responseData = await response.json();

		if (response.ok) {
			document.cookie = `jwt=${responseData.token}; path=/; secure;`;
			return true;
		} else {
			setWrongCredentials(true);
			return false;
		}
	} catch (error) {
		toast.error("Login Failed. Try again.");
		return false;
	}
};
