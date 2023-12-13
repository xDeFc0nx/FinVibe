/* eslint-disable jsx-a11y/control-has-associated-label */

"use client";

import React, { useState } from "react";

const page = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  return (
    <div>
      <div>
        <button
          type="button"
          onClick={() => {
            setSidebarOpen(!sidebarOpen);
          }}
          className="p-2 mr-3 text-gray-600 rounded cursor-pointer  hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-700"
        >
          <svg
            className="w-5 h-5"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 16 12"
          >
            <path
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M1 1h14M1 6h14M1 11h7"
            />
          </svg>
        </button>
        <div className={` ${sidebarOpen ? "flex  " : "hidden"}`}>
          <h1>testaaaaa</h1>
        </div>
      </div>
    </div>
  );
};

export default page;
