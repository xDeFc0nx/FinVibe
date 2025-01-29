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
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import * as z from "zod";
import { useUserData } from "../context/userData";
import { useWebSocket } from "../WebSocketProvidor";
import { saveAccount } from "../sidebar/account-switcher";
import { useEffect } from "react";
const formSchema = z.object({
	Type: z.string().min(1, "Account type is required"),
});

export default function CreateAccount() {
	const navigate = useNavigate();
	const { setAccounts, setActiveAccount, activeAccount } = useUserData();

	const { socket, isReady } = useWebSocket();

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
	});

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
						setAccounts((prevAccounts) => [...prevAccounts, response.account]);
						setActiveAccount(response.account);
						//navigate("/app/dashboard")
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
			console.log(activeAccount);

			localStorage.setItem("activeAccount", JSON.stringify(activeAccount));
			console.log("Saved Account to LocalStorage:", activeAccount);
		}
	}, [activeAccount]);
	return (
		<>
			<div className="container relative hidden h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-1 lg:px-0">
				<Link to="/app/dashboard">
					<Button className="absolute right-4 top-4 md:right-8 md:top-8">
						Login
					</Button>
				</Link>

				<div className="lg:p-8">
					<div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
						<div className="flex flex-col space-y-2 text-center">
							<h1 className="text-2xl font-semibold tracking-tight">
								{" "}
								Welcome!
							</h1>
							<p className="text-sm text-muted-foreground">
								Get ready for a smooth ride to financial freedom—no boring
								spreadsheets here. Let’s make money management a little more
								fun, shall we?
							</p>
						</div>

						<Form {...form}>
							<form onSubmit={form.handleSubmit(handleSubmit)}>
								<FormField
									control={form.control}
									name="Type"
									render={({ field }) => (
										<FormItem>
											<FormLabel>First Name</FormLabel>
											<FormControl>
												<Input
													placeholder="Your First Name"
													type=""
													{...field}
												/>
											</FormControl>

											<FormMessage />
										</FormItem>
									)}
								/>

								<Button type="submit" className="mt-5">
									Continue
								</Button>
							</form>
						</Form>
					</div>
				</div>
			</div>
		</>
	);
}
