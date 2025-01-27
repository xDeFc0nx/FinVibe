'use client';
import { TrendingUp } from 'lucide-react';
import { CartesianGrid, Line, LineChart, XAxis } from 'recharts';
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart';
import { useUserData } from '../context/userData';

const chartConfig = {
  Income: {
    label: 'Income',
    color: 'hsl(var(--chart-1))',
  },
  Expenses: {
    label: 'Expenses',
    color: 'hsl(var(--chart-5))',
  },
} satisfies ChartConfig;

export function BalanceChart() {
  const { activeAccount } = useUserData();

 const {chartOverview} = useUserData()

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
  const { activeAccount } = useUserData();

 const {chartOverview} = useUserData()

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
  const { activeAccount } = useUserData();

 const {chartOverview} = useUserData()

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
