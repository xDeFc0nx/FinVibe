import {
	Breadcrumb,
	BreadcrumbItem,
	BreadcrumbLink,
	BreadcrumbList,
	BreadcrumbPage,
	BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import { SidebarInset, SidebarTrigger } from "@/components/ui/sidebar";
import { Link, useNavigate } from "react-router";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
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
import { toast } from "react-toastify";
import { useWebSocket } from "@/components/WebSocketProvidor";
import { useUserData } from "@/components/context/userData";
import { ThemeChanger } from "@/components/ui/theme";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";

const formSchema = z.object({
	FirstName: z.string(),
	LastName: z.string(),
	Email: z.string(),
	Country: z.string(),
	OldPassword: z.string().min(8).optional(),
	NewPassword: z.string().min(8).optional(),
	ConfirmPassword: z.string().min(8).optional(),
});
export default function Index() {
	const { socket, isReady } = useWebSocket();

	const navigate = useNavigate();

	const { userData, setUserData } = useUserData();

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
	});
	function onSubmit(values: z.infer<typeof formSchema>) {
		try {
			console.log(values);
			if (socket && isReady) {
				socket.send("updateUser", {
					FirstName: values.FirstName,
					LastName: values.LastName,
					Email: values.Email,
					Country: values.Country,
					ID: userData?.ID,
				});

				socket.onMessage((msg) => {
					const response = JSON.parse(msg);

					if (response.userData) {
						setUserData({
							...userData,
							FirstName: values.FirstName,
							LastName: values.LastName,
							Email: values.Email,
							Country: values.Country,
						});
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
		if (userData) {
			form.setValue("FirstName", userData.FirstName);
			form.setValue("LastName", userData.LastName);
			form.setValue("Email", userData.Email);
			form.setValue("Country", userData.Country);
		}
	}, [userData, form]);

	const handleDelete = async () => {
		await fetch("http://localhost:3001/logout", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			credentials: "include",
		});
		navigate("/register");

		try {
			if (socket && isReady) {
				socket.send("deleteUser");

				socket.onMessage((msg) => {
					const response = JSON.parse(msg);

					if (response.Success) {
						toast.success("Account Deleted!");
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
	};

	return (
		<SidebarInset>
			<header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
				<div className="flex items-center gap-2 px-4">
					<SidebarTrigger className="-ml-1" />
					<Separator orientation="vertical" className="mr-2 h-4" />
					<Breadcrumb>
						<BreadcrumbList>
							<BreadcrumbItem className="hidden md:block">
								<Link to="/app/dashboard">
									<BreadcrumbLink>Dashboard</BreadcrumbLink>
								</Link>
							</BreadcrumbItem>
							<BreadcrumbSeparator className="hidden md:block" />
							<BreadcrumbItem>
								<BreadcrumbPage>User Settings</BreadcrumbPage>
							</BreadcrumbItem>
						</BreadcrumbList>
					</Breadcrumb>
				</div>
			</header>
			<div className="flex flex-1 flex-col gap-4 p-4 pt-0">
				<Form {...form}>
					<form onSubmit={form.handleSubmit(onSubmit)} className="">
						<FormField
							control={form.control}
							name="FirstName"
							render={({ field }) => (
								<FormItem>
									<FormLabel>First Name</FormLabel>
									<FormControl>
										<Input placeholder="First Name" type="text" {...field} />
									</FormControl>

									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="LastName"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Last Name</FormLabel>
									<FormControl>
										<Input placeholder="Last Name" type="text" {...field} />
									</FormControl>

									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="Email"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Email</FormLabel>
									<FormControl>
										<Input placeholder="Email" type="email" {...field} />
									</FormControl>

									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="Country"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Country</FormLabel>
									<FormControl>
										<Input
											placeholder="Your Country this decided your Currency"
											type="text"
											{...field}
										/>
									</FormControl>

									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="OldPassword"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Old Password</FormLabel>
									<FormControl>
										<PasswordInput placeholder="Old Password" {...field} />
									</FormControl>

									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="NewPassword"
							render={({ field }) => (
								<FormItem>
									<FormLabel>New Password</FormLabel>
									<FormControl>
										<PasswordInput placeholder="New Password" {...field} />
									</FormControl>

									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="ConfirmPassword"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Confirm Password</FormLabel>
									<FormControl>
										<PasswordInput placeholder="Confirm Password" {...field} />
									</FormControl>

									<FormMessage />
								</FormItem>
							)}
						/>

						<Button type="submit" variant="green" className="mt-5">
							Update
						</Button>
						<h3 className="mt-5">Theme</h3>
						<ThemeChanger />
					</form>
				</Form>
				<Dialog>
					<DialogTrigger>
						<Button variant="destructive">Delete Account</Button>
					</DialogTrigger>

					<DialogContent>
						<DialogHeader>
							<DialogTitle>Are you absolutely sure?</DialogTitle>
							<DialogDescription>
								This action cannot be undone. This will permanently delete your
								account and remove your data from our servers. <br />
								<Button
									variant="destructive"
									className="mt-5"
									onClick={handleDelete}
								>
									Im Sure!
								</Button>
							</DialogDescription>
						</DialogHeader>
					</DialogContent>
				</Dialog>
			</div>
		</SidebarInset>
	);
}
