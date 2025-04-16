"use client";

import * as React from "react"; // Import React for JSX and useMemo
import type { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal } from "lucide-react";

// UI Components
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator, // Keep if needed later
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

// Hooks & Context
import { useWebSocket } from "../WebSocketProvidor"; // Adjust path if needed
import { useSelector, useDispatch } from 'react-redux'; // Import hooks

// State & Actions & Types
import type { RootState, AppDispatch } from '@/store/store.ts'; // Adjust path
import { removeTransaction } from "@/store/slices/transactionsSlice"; // Action for deleting
// Import other actions/types if needed by backend responses
import { accountsReceived } from '@/store/slices/accountsSlice';
import { setChartData, setIncomePieData, setExpensePieData } from '@/store/slices/overviewSlice';
import type { Transaction, Account, UserData } from "@/types"; // Adjust path

// Notifications
import { toast } from "sonner"; // Or import from 'react-toastify'

// --- NO top-level hook calls here ---

export const columns: ColumnDef<Transaction>[] = [
  // --- Select Column (No Hooks Needed) ---
  {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && "indeterminate")
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
  // --- CreatedAt Column (No Hooks Needed) ---
  {
    accessorKey: "CreatedAt", // Verify casing matches your Transaction type
    header: "Created At",
    cell: ({ row }) => {
      const rawDate = row.getValue("CreatedAt") as string;
      try {
        const date = new Date(rawDate);
        if (isNaN(date.getTime())) return <div>Invalid Date</div>;
        return (
          <div>
            {date.toLocaleString("en-US", { // Consider user locale
              dateStyle: "short",
              timeStyle: "short",
            })}
          </div>
        );
      } catch(e) {
        console.error("Date parse error:", e);
        return <div>Error</div>
      }
    },
  },
  // --- Description Column (No Hooks Needed) ---
  {
    accessorKey: "Description", // Verify casing matches your Transaction type
    header: "Description",
  },
  // --- Type Column (No Hooks Needed) ---
  {
    accessorKey: "Type", // Verify casing matches your Transaction type
    header: () => <div className="text-right">Type</div>,
    cell: ({ row }) => <div className="text-right">{row.original.Type}</div>,
  },
  // --- Amount Column (Needs Hook -> Inline Component) ---
  {
    accessorKey: "Amount", // Verify casing matches your Transaction type
    header: () => <div className="text-right">Amount</div>,
    cell: ({ row }) => { // Outer cell function

      // --- Define Inline Component for this cell ---
      // Note: Defined newly on each cell render - less performant
      const InlineAmountCell = () => {
        // ✅✅✅ HOOKS ARE VALID HERE ✅✅✅
        const userData = useSelector((state: RootState) => state.user.data);

        // --- Logic using row data and hook data ---
        const amount = row.original.Amount;
        const type = row.original.Type;
        const currency = userData?.Currency ?? "$";

        const numericAmount = typeof amount === 'number' ? amount : parseFloat(amount as any);
        const formattedAmount = !isNaN(numericAmount)
          ? `${type === "Income" ? '+' : '-'}${currency}${numericAmount.toFixed(2)}`
          : `${currency} N/A`;

        const backgroundColor =
          type === "Income" ? "bg-green-400/50" : "bg-red-400/50";

        return (
          <div className={`text-right p-2 rounded-lg ${backgroundColor}`}>
            {formattedAmount}
          </div>
        );
      }; // End InlineAmountCell definition

      // --- Render the inline component ---
      return <InlineAmountCell />;

    }, // End cell function for Amount
  }, // End Amount column definition
  // --- Actions Column (Needs Hooks -> Inline Component) ---
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => { // Outer cell function

      // --- Define Inline Component for this cell ---
      // Note: Defined newly on each cell render - less performant
      const InlineActionsCell = () => {
        // ✅✅✅ HOOKS ARE VALID HERE ✅✅✅
        const dispatch: AppDispatch = useDispatch();
        const { socket, isReady } = useWebSocket();
        // Only select what's needed for delete logic (activeAccountId for check)
        const activeAccountId = useSelector((state: RootState) => state.accounts.activeAccountId);
        const currentAccounts = useSelector((state: RootState) => state.accounts.list); // Needed if updating totals
        const transaction = row.original;

        // --- Logic using hooks and row data ---
        const handleDelete = () => {
          if (!socket || !isReady) { toast.error("Connection not ready."); return; }
          if (!activeAccountId) { toast.error("No active account selected."); return; } // Check if needed? Delete might be global.
          if (!transaction?.ID) { toast.error("Invalid transaction ID."); return; } // Use uppercase ID if that's correct

          const transactionIdToDelete = transaction.ID; // Use correct ID

          // Temporary handler
          const handleDeleteResponse = (msg: string) => {
            let shouldRemoveListener = true; // Assume we remove unless specified otherwise
            try {
              const response = JSON.parse(msg);
              console.log("Delete Response:", response);

              // Check if this is the relevant success message
              // Backend should ideally send correlation ID or specific action type
              if (response.Success && response.deletedId === transactionIdToDelete) {
                dispatch(removeTransaction(transactionIdToDelete));
                toast.success("Transaction deleted!");

                // Optional: Update account/overview state if backend provides data
                if (response.AccountData && activeAccountId) {
                  dispatch(updateAccountDetails({ id: activeAccountId, details: response.AccountData }));
                }
                if (response.OverviewData) {
                  // dispatch appropriate overview action(s)
                }

              } else if (response.Error /* && check if related error */) {
                toast.error(`Failed to delete: ${response.Error}`);
              } else {
                // This message wasn't for us, don't remove the listener yet
                shouldRemoveListener = false;
              }

            } catch (parseError) {
              console.error("Error parsing delete response:", parseError, msg);
              // Still remove listener on parse error
            } finally {
              if (shouldRemoveListener) {
                 console.log("Removing handleDeleteResponse listener");
                 // Check if socket still exists before calling offMessage
                 if (socket) {
                    socket.offMessage(handleDeleteResponse); // Assumes offMessage exists
                 }
              }
            }
          }; // End handleDeleteResponse

          // Send request
          try {
            console.log("Registering delete listener for:", transactionIdToDelete);
            socket.onMessage(handleDeleteResponse);
            socket.send("deleteTransaction", { ID: transactionIdToDelete }); // Use correct ID
            console.log("Sent deleteTransaction request for ID:", transactionIdToDelete);
          } catch (sendError) {
            console.error("Failed to send delete request:", sendError);
            toast.error("Failed to send delete request.");
            // Clean up listener if send fails
             if (socket) {
                socket.offMessage(handleDeleteResponse);
             }
          }
        }; // End handleDelete

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem
                onClick={handleDelete}
                className="text-red-600 focus:text-red-700 focus:bg-red-100"
              >
                Delete
              </DropdownMenuItem>
              {/* <DropdownMenuSeparator /> */}
            </DropdownMenuContent>
          </DropdownMenu>
        );
      }; // End InlineActionsCell definition

      // --- Render the inline component ---
      return <InlineActionsCell />;

    }, // End cell function for Actions
  }, // End Actions column definition
]; // End columns array
