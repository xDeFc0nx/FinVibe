"use client";
import { Line, LineChart, Pie, PieChart, Cell } from "recharts";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../ui/card";

import { useSelector } from 'react-redux';
import type { RootState } from '@/store/store.ts';
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
  const chartData = useSelector((state: RootState) => state.overview.chartData)
  return (
    <ChartContainer config={chartConfig}>
      <LineChart
        accessibilityLayer
        data={chartData}
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

  const chartData = useSelector((state: RootState) => state.overview.chartData)
  return (
    <ChartContainer config={chartConfig}>
      <LineChart
        accessibilityLayer
        data={chartData}
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

  const chartData = useSelector((state: RootState) => state.overview.chartData)
  return (
    <ChartContainer config={chartConfig}>
      <LineChart
        accessibilityLayer
        data={chartData}
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
  const incomePieData = useSelector((state: RootState) => state.overview.incomePieData);
  const colorPalette = [
    "hsl(var(--chart-1))",
    "hsl(var(--chart-2))",
    "hsl(var(--chart-3))",
    "hsl(var(--chart-4))",
    "hsl(var(--chart-5))",
    "hsl(var(--chart-6))",
    "hsl(var(--chart-7))",
    "hsl(var(--chart-8))",
    "hsl(var(--chart-9))",
    "hsl(var(--chart-10))",
  ];

  return (
    <Card className="flex flex-col">
      <CardContent className="flex-1 pb-0">
        <CardHeader className="items-center pb-0">
          <CardTitle>Income</CardTitle>
        </CardHeader>

        <ChartContainer config={chartConfig} className="mx-auto aspect-square ">
          <PieChart>
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />

            <Pie
              data={incomePieData}
              dataKey="Amount"
              nameKey="Description"
              fill="var(--color-Income)"
            >
              {incomePieData.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={colorPalette[index % colorPalette.length]}
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

  const expensesPieData = useSelector((state: RootState) => state.overview.expensePieData);
  const colorPalette = [
    "hsl(var(--chart-1))",
    "hsl(var(--chart-2))",
    "hsl(var(--chart-3))",
    "hsl(var(--chart-4))",
    "hsl(var(--chart-5))",
    "hsl(var(--chart-6))",
    "hsl(var(--chart-7))",
    "hsl(var(--chart-8))",
    "hsl(var(--chart-9))",
    "hsl(var(--chart-10))",
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
              data={expensesPieData}
              dataKey="Amount"
              nameKey="Description"
              fill="var(--color-Expenses)"
            >
              {expensesPieData.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={colorPalette[index % colorPalette.length]}
                />
              ))}
            </Pie>
          </PieChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
