import React from "react";

import { motion } from "framer-motion";
import Image from "next/image";
import Illustration from "./illustration";
import Github from "../../Icon/github.svg";

const hero = () => (
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
          className="bg-primary-blue   justify-center gap-1 w-32 h-10 font-bold p-2 rounded-full text-white  "
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
          className=" flex justify-center gap-1 border border-[#3F3F46] font-bold p-2 w-32 h-10 rounded-full text-white "
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
);

export default hero;
