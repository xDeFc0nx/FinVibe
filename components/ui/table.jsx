/* eslint-disable jsx-a11y/control-has-associated-label */

"use client";

import React, { useEffect, useState } from "react";
import { toast, ToastContainer } from "react-toastify";

import NewTransaction from "./newTransaction";

export default function Tables() {
  const [transactions, setTransactions] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await fetch("http://localhost:3000/api/transactions", {
          method: "PUT",
        });

        if (!res.ok) {
          throw new Error("Failed to fetch data");
        }

        const data = await res.json();
        setTransactions(data);
      } catch (error) {
        console.error(error);
        toast.error("Failed to fetch data", {
          position: "bottom-right",
        });
      }
    };

    fetchData();
  }, []);

  async function deleteData(id) {
    try {
      const pos = await fetch(`http://localhost:3000/api/transactions/${id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id,
        }),
      });

      if (!pos.ok) {
        throw new Error();
      } else {
        toast.success("Deleted Data", {
          position: "bottom-right",
        });
      }

      // Handle response if needed
      const result = await pos.json();
      console.log(result);
    } catch (error) {
      console.error(error);
      toast.error("Failed to delete data", {
        position: "bottom-right",
      });
    }
  }

  return (
    <div className="relative overflow-x-auto shadow-md sm:rounded-lg top-20">
      <NewTransaction />
      <ToastContainer />
      <table className="w-[95%] text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
        <thead className="text-xs text-700 uppercase dark:bg-[#24303F]/50 dark:text-gray-400">
          <tr>
            <th scope="col" className="px-6 py-3">
              Transaction
            </th>
            <th scope="col" className="px-6 py-3">
              Amount
            </th>
            <th scope="col" className="px-6 py-3">
              Description
            </th>
            <th scope="col" className="px-6 py-3">
              Date
            </th>
            <th scope="col" className="px-6 py-3">
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          {transactions &&
            transactions.map((transaction) => (
              <tr
                key={transaction.id}
                className="bg-white border-b dark:bg-[#24303F]/50 dark:border-gray-700"
              >
                <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                  {transaction.type}
                </td>
                <td className="px-6 py-4">${transaction.amount}</td>
                <td className="px-6 py-4">{transaction.description}</td>
                <td className="px-6 py-4">{transaction.DateCreated}</td>

                <td className="px-6 py-4 text-right ">
                  <button
                    href={transaction.id}
                    type="submit"
                    className="text-blue-500 px-3"
                  >
                    Edit
                  </button>

                  <button
                    onClick={() => deleteData(transaction.id)}
                    type="submit"
                    className="text-red-500"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
        </tbody>
      </table>
    </div>
  );
}
