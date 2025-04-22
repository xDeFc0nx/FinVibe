"use client";
import type { ColumnDef } from "@tanstack/react-table";
import type { RootState, } from '@/store/store.ts';
import type { Transaction, } from "@/types";
import { useSelector, } from 'react-redux';
export const columns: ColumnDef<Transaction>[] = [
  {
    accessorKey: "Description",
    header: "Description",
  },
  {
    accessorKey: "type",
    header: () => <div className="text-right">Type</div>,
    cell: ({ row }) => <div className="text-right">{row.original.Type}</div>,
  },

  {
    accessorKey: "amount",
    header: () => <div className="text-right">Amount</div>,
    cell: ({ row }) => {
      const amount = row.original.Amount;
      const type = row.original.Type;
      const UserData = useSelector((state: RootState) => state.user.data)

      const formattedAmount =
        type === "Income"
          ? `+${UserData?.Currency}${amount}`
          : `-${UserData?.Currency}${amount}`;

      const backgroundColor =
        type === "Income" ? "bg-green-400/50" : "bg-red-400/50";

      return (
        <div className={`text-right p-2 rounded-lg ${backgroundColor}`}>
          {formattedAmount}
        </div>
      );
    },
  },
];
