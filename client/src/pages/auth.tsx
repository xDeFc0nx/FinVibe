"use client";
import { Card, CardHeader } from "@/components/ui/card";
import { useState } from "react";
import Register from "@/components/ui/register";
import Login from "@/components/ui/login";

export default function Auth() {
	const [activeTab, setActiveTab] = useState<"register" | "login">("register");

	return (
		<Card className="mx-auto mt-20  flex w-full flex-col justify-center  sm:w-[500px]">
			<div className="inline w-full mt-5">
				<button
					onClick={() => setActiveTab("register")}
					className="w-[50%] text-lg relative"
				>
					Register
					<div
						className={`h-1 transition-colors ${
							activeTab === "register" ? "bg-blue-500" : "bg-gray-500"
						}`}
					/>
				</button>

				<button
					onClick={() => setActiveTab("login")}
					className="w-[50%] text-lg relative"
				>
					Login
					<div
						className={`h-1 transition-colors ${
							activeTab === "login" ? "bg-blue-500" : "bg-gray-500"
						}`}
					/>
				</button>
			</div>
			<CardHeader></CardHeader>

			{activeTab === "register" ? <Register /> : <Login />}
		</Card>
	);
}
