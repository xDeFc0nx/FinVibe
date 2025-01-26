"use client"
import { TrendingUp } from "lucide-react"
import { CartesianGrid, Line, LineChart, XAxis } from "recharts"
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"
import { useUserData } from "../context/userData"

const chartConfig = {
  desktop: {
    label: "Income",
    color: "hsl(var(--chart-1))",
  },
  mobile: {
    label: "Expenses",
    color: "hsl(var(--chart-2))",
  },
} satisfies ChartConfig

export function BalanceChart() {
  const { activeAccount } = useUserData();

  const chartData = [
    { month: "January", Income: activeAccount?.Income, Expenses: activeAccount?.Expense },
    { month: "February", Income: activeAccount?.Income, Expenses: activeAccount?.Expense },
    { month: "March", Income: activeAccount?.Income, Expenses: activeAccount?.Expense },
    { month: "April", Income: activeAccount?.Income, Expenses: activeAccount?.Expense },
    { month: "May", Income: activeAccount?.Income, Expenses: activeAccount?.Expense },
    { month: "June", Income: activeAccount?.Income, Expenses: activeAccount?.Expense },
  ]

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
          stroke="var(--color-desktop)"
          strokeWidth={2}
          dot={false}
        />
        <Line
          dataKey="Expenses"
          type="monotone"
          stroke="var(--color-mobile)"
          strokeWidth={2}
          dot={false}
        />
      </LineChart>
    </ChartContainer>
  )
}
