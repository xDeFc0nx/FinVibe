"use client"

import {type ColumnDef } from "@tanstack/react-table"
import { useUserData, type Transaction } from "@/components/context/userData"; 

export const columns: ColumnDef<Transaction>[] = [
  {
    accessorKey: "CreatedAt",
    header: "Created At",
  },
  {
    accessorKey: "Description",
    header: "Description",
  },
  {
    accessorKey: "Amount",
    header: "Amount",
  },
]

