/* eslint-disable jsx-a11y/control-has-associated-label */

"use client";

import React from "react";
import NewTransaction from "./newTransaction";

async function getData() {
  const res = await fetch("http://localhost:3000/api/transactions", {
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error("Failed to fetch data");
  }
  return res.json();
}

export default async function Tables() {
  const transactions = await getData();

  return (
    <div className="relative overflow-x-auto shadow-md sm:rounded-lg top-20">
      <NewTransaction />
      <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
        <thead className="text-xs w-full text-700 uppercase bg-gray-[#24303F] dark:bg-[#24303F] dark:text-gray-400">
          <tr>
            <th scope="col" className="px-6 py-3">
              Transaction
            </th>
            <th scope="col" className="px-6 py-3">
              <div className="flex items-center">Amount</div>
            </th>
            <th scope="col" className="px-6 py-3">
              <div className="flex items-center">Description</div>
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
                className="bg-white border-b dark:bg-gray-800 dark:border-gray-700"
              >
                <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                  {transaction.type}
                </td>
                <td className="px-6 py-4">${transaction.amount}</td>
                <td className="px-6 py-4">{transaction.description}</td>
                <td className="px-6 py-4 text-right">
                  <button
                    href={transaction.id}
                    type="submit"
                    className="text-blue-500"
                  >
                    Edit
                  </button>
                </td>
              </tr>
            ))}
        </tbody>
      </table>
    </div>
  );
}
