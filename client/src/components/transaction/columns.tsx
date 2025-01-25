'use client';
import type { ColumnDef } from '@tanstack/react-table';
import { MoreHorizontal } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

import { useUserData, type Transaction } from '@/components/context/userData';

export const columns: ColumnDef<Transaction>[] = [
  {
    id: 'select',
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && 'indeterminate')
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'CreatedAt',
    header: 'Created At',
    cell: ({ row }) => {
      const rawDate = row.getValue('CreatedAt') as string;
      const date = new Date(rawDate);

      if (isNaN(date.getTime())) throw new Error('Invalid date');

      return (
        <div>
          {date.toLocaleString('en-US', {
            dateStyle: 'short',
            timeStyle: 'medium',
          })}
        </div>
      );
    },
  },
  {
    accessorKey: 'Description',
    header: 'Description',
  },
  {
    accessorKey: 'type',
    header: () => <div className="text-right">Type</div>,
    cell: ({ row }) => <div className="text-right">{row.original.Type}</div>,
  },

 {
  accessorKey: 'amount',
  header: () => <div className="text-right">Amount</div>,
  cell: ({ row }) => {
    const amount = row.original.Amount;
    const type = row.original.Type;
    
    // Add + for Income, - for Expense
    const formattedAmount = type === 'Income' ? `+${amount}` : `-${amount}`;
    
    const backgroundColor = type === 'Income' ? 'bg-green-400/50' : 'bg-red-400/50';
    
    return (
      <div className={`text-right p-2 rounded-lg ${backgroundColor}`}>
        {formattedAmount}
      </div>
    );
  },
},  {
    id: 'actions',
    enableHiding: false,
    cell: ({ row }) => {
      const payment = row.original;

      const transaction = row.original;
      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(transaction.ID)}
            >
              Copy payment ID
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>View customer</DropdownMenuItem>
            <DropdownMenuItem>View payment details</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];
