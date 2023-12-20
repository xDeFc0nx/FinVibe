/* eslint-disable jsx-a11y/control-has-associated-label */

"use client";

import React, { useEffect, useState } from "react";
import { toast, ToastContainer } from "react-toastify";
import { FetchData } from "../../../actions/fetch";
import { DeleteData } from "../../../actions/Delete";

const Tables = () => {
  const [transactions, setTransactions] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await FetchData();
        setTransactions(data);
      } catch (error) {
        console.error(error);
        toast.error(error.message);
      }
    };

    fetchData();
  }, []);
  async function deleteData(id) {
    try {
      const Delete = await DeleteData({ params: { id } });

      if (Delete) {
        toast.success("Deleted Data", {
          position: "bottom-right",
        });

        // Handle response if needed
        console.log(Delete);
      } else {
        throw new Error("Invalid response from server");
      }
    } catch (error) {
      console.error(error);
      toast.error("Failed to delete data", {
        position: "bottom-right",
      });
    }
  }

  return (
    <div className=" overflow-x-auto shadow-md sm:rounded-lg top-20 ">
      <ToastContainer />
      <table className="w-full text-sm text-left rtl:text-right text-gray-500  bg-secondary-gray/50  backdrop-filter backdrop-blur-lg	    shadow-lg dark:text-gray-400">
        <thead className="text-xs text-700 uppercase bg-secondary-gray/50  backdrop-filter backdrop-blur-lg	    shadow-lg dark:text-gray-400">
          <tr>
            <th scope="col" className="px-6 py-3">
              Date
            </th>

            <th scope="col" className="px-6 py-3">
              Description
            </th>
            <th scope="col" className="px-6 py-3">
              Amount
            </th>
            <th scope="col" className="px-6 py-3">
              Category
            </th>
          </tr>
        </thead>
        <tbody>
          {transactions &&
            transactions.map((transaction) => (
              <tr
                key={transaction.id}
                className="border-b border-black bg-secondary-gray/50  backdrop-filter backdrop-blur-lg	    shadow-lg"
              >
                <td className="px-6 py-4">
                  {new Date(transaction.DateCreated).toLocaleDateString()}
                </td>
                <td className="px-6 py-4">{transaction.description}</td>
                <td className="px-6 py-4">${transaction.amount}</td>
                <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                  {transaction.type}
                </td>

                <td>
                  <button
                    href={transaction.id}
                    type="submit"
                    className="text-blue-500 px-3"
                  >
                    Edit
                  </button>
                </td>
                <td>
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
};
export default Tables;
