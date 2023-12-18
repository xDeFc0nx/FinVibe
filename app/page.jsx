"use client";

import React from "react";
import { motion } from "framer-motion";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";
import Image from "next/image";
import Card from "../components/ui/Card";
import Bg from "../components/ui/LandingPage/background";
import Illustration from "../components/ui/LandingPage/illustration";
import RealTimeSvg from "../components/Icon/RealTimeSvg.svg";
import PersonalizedSvg from "../components/Icon/PersonalizedSvg.svg";
import SecureSvg from "../components/Icon/SecureSvg.svg";
import ExpenseSvg from "../components/Icon/ExpenseSvg.svg";
import Github from "../components/Icon/github.svg";

const page = () => (
  <>
    <ToastContainer />
    <Bg />
    <div className="flex ">
      <div className="hero pt-32 pl-32">
        <p className="text-5xl text-white font-bold">
          FinVibe <span className="text-primary-pink"> Elevate </span>
          Your Financial
          <br /> Journey with Empowering Tools
          <br /> and Insightful Guidance
        </p>

        <p className="text-lg text-[#95959D]  pt-4 pb-4">
          Empower your finances with FinVibe's smart tools and expert guidance
        </p>
        <div className="space-x-5 flex">
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            type="button"
            className="bg-primary-blue w-32 h-10 font-bold p-2 rounded-full text-white  "
          >
            Get Started{" "}
            <svg
              className="inline"
              width="20"
              height="10"
              viewBox="0 0 20 10"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M15.0125 3.75H0V6.25H15.0125V10L20 5L15.0125 0V3.75Z"
                fill="white"
              />
            </svg>
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            type="button"
            className=" flex gap-1 border border-[#3F3F46] font-bold p-2 w-32 h-10 rounded-full text-white "
          >
            <Image src={Github} width={20} height={20} />
            Github
          </motion.button>
        </div>
      </div>
      <div className="pl-56">
        <Illustration />
      </div>
    </div>
    <div className=" pt-32 pl-32 flex space-x-14">
      <Card
        icon={RealTimeSvg}
        header="Real-time Financial Insights"
        value="Display real-time updates on financial data, providing users with the latest information on their accounts, transactions, and investment performance."
      />
      <Card
        icon={PersonalizedSvg}
        header="Personalized Budgeting"
        value="Ensure the security of financial transactions by offering a feature that monitors and alerts users about potentially suspicious or unauthorized activities."
      />
      <Card
        icon={SecureSvg}
        header="Secure Transaction Monitoring"
        value="Provide users with detailed analytics on their investments, including performance charts, historical data, and recommendations for optimizing their investment portfolio."
      />
      <Card
        icon={ExpenseSvg}
        header="Expense Categorization"
        value="Automatically categorize expenses to help users understand their spending habits better. This feature can simplify budgeting and enable users to identify areas where they can save money."
      />
    </div>
  </>
);

export default page;
