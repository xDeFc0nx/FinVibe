"use client";
import { TrendingUp } from "lucide-react";
import { Line, LineChart, Label, Pie, PieChart, Cell } from "recharts";
import {
	type ChartConfig,
	ChartContainer,
	ChartTooltip,
	ChartTooltipContent,
} from "@/components/ui/chart";
import { useUserData } from "../context/userData";
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from "../ui/card";

const chartConfig = {
	Income: {
		label: "Income",
		color: "hsl(var(--chart-1))",
	},
	Expenses: {
		label: "Expenses",
		color: "hsl(var(--chart-5))",
	},
} satisfies ChartConfig;

export function BalanceChart() {
	const { chartOverview } = useUserData();

	return (
		<ChartContainer config={chartConfig}>
			<LineChart
				accessibilityLayer
				data={chartOverview}
				margin={{
					left: 12,
					right: 12,
				}}
			>
				<ChartTooltip cursor={false} content={<ChartTooltipContent />} />
				<Line
					dataKey="Income"
					type="monotone"
					stroke="var(--color-Income)"
					strokeWidth={2}
					dot={false}
				/>
				<Line
					dataKey="Expenses"
					type="monotone"
					stroke="var(--color-Expenses)"
					strokeWidth={2}
					dot={false}
				/>
			</LineChart>
		</ChartContainer>
	);
}
export function IncomeChart() {
	const { chartOverview } = useUserData();

	return (
		<ChartContainer config={chartConfig}>
			<LineChart
				accessibilityLayer
				data={chartOverview}
				margin={{
					left: 12,
					right: 12,
				}}
			>
				<ChartTooltip cursor={false} content={<ChartTooltipContent />} />
				<Line
					dataKey="Income"
					type="monotone"
					stroke="var(--color-Income)"
					strokeWidth={2}
					dot={false}
				/>
			</LineChart>
		</ChartContainer>
	);
}
export function ExpensesChart() {
	const { chartOverview } = useUserData();

	return (
		<ChartContainer config={chartConfig}>
			<LineChart
				accessibilityLayer
				data={chartOverview}
				margin={{
					left: 12,
					right: 12,
				}}
			>
				<ChartTooltip cursor={false} content={<ChartTooltipContent />} />
				<Line
					dataKey="Expenses"
					type="monotone"
					stroke="var(--color-Expenses)"
					strokeWidth={2}
					dot={false}
				/>
			</LineChart>
		</ChartContainer>
	);
}
export function IncomePie() {
	const { incomePie } = useUserData(); // Use incomePie data from context

	// Define a color palette for the pie chart slices
	const colorPalette = [
		"hsl(var(--chart-1))", // Color 1
		"hsl(var(--chart-2))", // Color 2
		"hsl(var(--chart-3))", // Color 3
		"hsl(var(--chart-4))", // Color 4
		"hsl(var(--chart-5))", // Color 5
		"hsl(var(--chart-6))", // Color 6
		"hsl(var(--chart-7))", // Color 7
		"hsl(var(--chart-8))", // Color 8
		"hsl(var(--chart-9))", // Color 9
		"hsl(var(--chart-10))", // Color 10
	];

	return (
		<Card className="flex flex-col">
			<CardContent className="flex-1 pb-0">
				<CardHeader className="items-center pb-0">
					<CardTitle>Expenses</CardTitle>
				</CardHeader>

				<ChartContainer config={chartConfig} className="mx-auto aspect-square ">
					<PieChart>
						<ChartTooltip
							cursor={false}
							content={<ChartTooltipContent hideLabel />}
						/>
						<ChartTooltip
							cursor={false}
							content={<ChartTooltipContent hideLabel />}
						/>
						<Pie
							data={incomePie} // Use incomePie data
							dataKey="Amount" // Key for the value in the data
							nameKey="Description" // Key for the label in the data
							fill="var(--color-Income)" // Default fill color
						>
							{incomePie.map((entry, index) => (
								<Cell
									key={`cell-${index}`}
									fill={colorPalette[index % colorPalette.length]} // Assign unique color
								/>
							))}
						</Pie>
					</PieChart>
				</ChartContainer>
			</CardContent>
		</Card>
	);
}
export function ExpensesPie() {
	const { expensesPie } = useUserData(); // Use expensesPie data from context

	// Define a color palette for the pie chart slices
	const colorPalette = [
		"hsl(var(--chart-1))", // Color 1
		"hsl(var(--chart-2))", // Color 2
		"hsl(var(--chart-3))", // Color 3
		"hsl(var(--chart-4))", // Color 4
		"hsl(var(--chart-5))", // Color 5
		"hsl(var(--chart-6))", // Color 6
		"hsl(var(--chart-7))", // Color 7
		"hsl(var(--chart-8))", // Color 8
		"hsl(var(--chart-9))", // Color 9
		"hsl(var(--chart-10))", // Color 10
	];

	return (
		<Card className="flex flex-col">
			<CardHeader className="items-center pb-0">
				<CardTitle>Expenses</CardTitle>
			</CardHeader>
			<CardContent className="flex-1 pb-0">
				<ChartContainer config={chartConfig} className="mx-auto aspect-square ">
					<PieChart>
						<ChartTooltip
							cursor={false}
							content={<ChartTooltipContent hideLabel />}
						/>
						<Pie
							data={expensesPie} // Use expensesPie data
							dataKey="Amount" // Key for the value in the data
							nameKey="Description" // Key for the label in the data
							fill="var(--color-Expenses)" // Default fill color
						>
							{expensesPie.map((entry, index) => (
								<Cell
									key={`cell-${index}`}
									fill={colorPalette[index % colorPalette.length]} // Assign unique color
								/>
							))}
						</Pie>
					</PieChart>
				</ChartContainer>
			</CardContent>
		</Card>
	);
}
