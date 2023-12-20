/* eslint-disable import/no-extraneous-dependencies */

"use client";

import React from "react";

import { motion } from "framer-motion";

export default function newTransaction({ color, Text, type }) {
  return (
    <div>
      <motion.button
        whileHover={{ scale: 1.1 }}
        whileTap={{ scale: 0.85 }}
        type={type}
        className={`py-2.5 px-5 me-2 mb-2 text-sm font-medium rounded-lg ${color} `}
      >
        {Text}
      </motion.button>
    </div>
  );
}
