import React from "react";

import { motion } from "framer-motion";
import Image from "next/image";
import Link from "next/link";
import Illustration from "./illustration";
import Github from "../../Icon/github.svg";

const hero = () => (
  <div className="lg:flex md:grid-cols-1  ">
    <div className=" pt-32 pl-32  ">
      <p className=" lg:text-5xl md:text-lg sm:text-xs text-white font-bold">
        FinVibe <span className="text-primary-pink"> elevate </span>
        your financial
        <br /> journey with empowering tools
        <br /> and insightful guidance
      </p>

      <p className=" lg:text-lg md:text-sm sm:text-xs text-[#95959D]  pt-4 pb-4">
        Empower your finances with FinVibe's smart tools and expert guidance
      </p>
      <div className="space-x-5 flex">
        <Link href="/dashboard">
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            type="button"
            className="bg-primary-blue   flex items-center  space-x-2 w-40 h-10 font-bold px-4 py-2  rounded-full text-white  "
          >
            <span>Get Started</span>
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
        </Link>
        <Link href="https://github.com/xDeFc0nx/FinVibe">
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            type="button"
            className=" flex items-center  border border-[#3F3F46]  justify-center gap-1 w-36 h-10 font-bold px-4 py-2 rounded-full text-white "
          >
            <Image src={Github} width={20} height={20} />
            Github
          </motion.button>
        </Link>
      </div>
    </div>
    <div className="pl-56 pt-20 flex items-center ">
      <div
        className={` flex text-white w-[30rem] h-[20rem] bg-white/10  backdrop-filter backdrop-blur-md rounded-lg 	   shadow-lg`}
      >
        <Illustration />
      </div>
    </div>
  </div>
);

export default hero;
