/* eslint-disable import/no-extraneous-dependencies */

"use client";

import React, { useState } from "react";

import { AnimatePresence } from "framer-motion";
import axios from "axios";
import { toast } from "react-toastify";
import Modal from "./modal";

const newTransaction = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const close = () => setModalOpen(false);
  const open = () => setModalOpen(true);
  const [inputs, setInputs] = useState({});
  const transactionData = {
    type: inputs.type,
    amount: parseFloat(inputs.amount), // Convert amount to a number
    description: inputs.description,
  };
  const handleSubmit = async (e) => {
    e.preventDefault();
    axios
      .post("http://localhost:3000/api/transactions", transactionData)
      .then((res) => {
        console.log(res);
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        setInputs({});
        setModalOpen(false);
        toast("Added Transaction!", {
          type: "success",
          position: "bottom-right",
        });
      });
  };
  const handleChange = (e) => {
    const { name } = e.target;
    const { value } = e.target;
    setInputs((prevState) => ({ ...prevState, [name]: value }));
  };

  return (
    <div>
      <button
        type="button"
        onClick={() => (modalOpen ? close() : open())}
        className="w-10 h-10 rounded-sm bg-blue-500"
      >
        New
      </button>
      <AnimatePresence mode="wait">
        {modalOpen && (
          <Modal
            modalOpen={modalOpen}
            handleClose={close}
            text="New Transaction"
          >
            <form className="p-4 md:p-5 " onSubmit={handleSubmit}>
              <div className="col-span-2 sm:col-span-1">
                <label
                  htmlFor="price"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Type
                </label>
                <input
                  name="type"
                  value={inputs.type || ""}
                  onChange={handleChange}
                  type="text"
                  className=" border-[2px] border-gray-300 text-gray-900 text-sm rounded-lg block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500 outline-none"
                  required
                />
              </div>
              <div className="col-span-2 sm:col-span-1">
                <label
                  htmlFor="price"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Amount
                </label>
                <input
                  name="amount"
                  value={inputs.amount || ""}
                  onChange={handleChange}
                  type="number"
                  className="border-[2px] border-gray-300 text-gray-900 text-sm rounded-lg block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500 outline-none"
                  placeholder="$2999"
                  required
                />
              </div>
              <div className="col-span-2">
                <label
                  htmlFor="description"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Description
                </label>
                <textarea
                  name="description"
                  value={inputs.description || ""}
                  onChange={handleChange}
                  id="description"
                  rows="4"
                  className="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 outline-none"
                  placeholder="Write description here"
                  required
                />
              </div>
              <button
                type="submit"
                className="text-white inline-flex items-center bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
              >
                <svg
                  className="me-1 -ms-1 w-5 h-5"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    fillRule="evenodd"
                    d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z"
                    clipRule="evenodd"
                  />
                </svg>
                Add
              </button>
            </form>
          </Modal>
        )}
      </AnimatePresence>
    </div>
  );
};

export default newTransaction;
