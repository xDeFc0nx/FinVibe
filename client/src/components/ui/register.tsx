"use client";
import { Button } from "@/components/ui/button";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { PasswordInput } from "@/components/ui/password-input";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import {  useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import * as z from "zod";
import { handleLogin } from "@/components/ui/login";
import CreateAccount from "@/components/ui/createAccount";
import { useState } from "react";
import { WebSocketProvider } from "@/components/WebSocketProvidor";
const formSchema = z.object({
	firstName: z.string(),
	lastName: z.string(),
	currency: z.string(),
	email: z.string(),
	password: z.string(),
});

export default function Register() {
	const [wrongCredentials, setWrongCredentials] = useState(false);
	const navigate = useNavigate();
	const [registered, setRegistered] = useState(false);
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
	});

	const handleSubmit = async (data: z.infer<typeof formSchema>) => {
		console.log(data);

		try {
			const response = await fetch("http://localhost:3001/Register", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(data),
				credentials: "include",
			});

			if (response.ok) {
				await handleLogin(
					{ email: data.email, password: data.password },
					setWrongCredentials,
				);
				setRegistered(true);
			} else {
				toast("Login Failed Try again");
			}
		} catch (error) {
			toast("Login Failed Try again");
		}
	};
	if (registered) {
		return (
			<WebSocketProvider>
					<CreateAccount />;
			</WebSocketProvider>
		);
	}
	return (
		<div className="p-8">
			<div className="flex flex-col space-y-2 text-center">
				<h1 className="text-2xl font-semibold tracking-tight"> Welcome!</h1>
				<p className="text-sm text-muted-foreground">
					Get ready for a smooth ride to financial freedom—no boring
					spreadsheets here. Let’s make money management a little more fun,
					shall we?
				</p>
			</div>
			<br />
			<Form {...form}>
				<form onSubmit={form.handleSubmit(handleSubmit)}>
					<FormField
						control={form.control}
						name="firstName"
						render={({ field }) => (
							<FormItem>
								<FormControl>
									<Input placeholder="Your First Name" type="" {...field} />
								</FormControl>

								<FormMessage />
							</FormItem>
						)}
					/>

					<br />
					<FormField
						control={form.control}
						name="lastName"
						render={({ field }) => (
							<FormItem>
								<FormControl>
									<Input placeholder="Your Last Name" type="" {...field} />
								</FormControl>

								<FormMessage />
							</FormItem>
						)}
					/>

					<br />
					<FormField
						control={form.control}
						name="currency"
						render={({ field }) => (
							<FormItem>
								<FormControl>
									<Input placeholder="Your currency" type="" {...field} />
								</FormControl>

								<FormMessage />
							</FormItem>
						)}
					/>

					<br />
					<FormField
						control={form.control}
						name="email"
						render={({ field }) => (
							<FormItem>
								<FormControl>
									<Input placeholder="Your Email" type="" {...field} />
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
									<PasswordInput placeholder="Your Password" {...field} />
								</FormControl>

								<FormMessage />
							</FormItem>
						)}
					/>
					<br />
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

					<Button type="submit" className="mt-5">
						Continue
					</Button>
				</form>
			</Form>
		</div>
	);
}
