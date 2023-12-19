/* eslint-disable jsx-a11y/control-has-associated-label */

"use client";

import { UserButton } from "@clerk/nextjs";
import { useState } from "react";
import { motion } from "framer-motion";
import Image from "next/image";
import Link from "next/link";
import Logo from "../../Icon/Logo.svg";
import Links from "./links";
import dashboardSvg from "../../Icon/dashboardSvg.svg";
import TransactionSvg from "../../Icon/TransactionSvg.svg";

export default function Sidebar() {
  const variants = {
    open: {
      transition: { staggerChildren: 0.07, delayChildren: 0.2, duration: 0.9 },
    },
    closed: {
      transition: {
        staggerChildren: 0.05,
        staggerDirection: -1,
        duration: 0.9,
      },
    },
  };

  const [hidden, setHidden] = useState(true);

  return (
    <motion.aside
      initial={{
        width: hidden ? "20rem" : "4rem",
        height: "100vh",
      }}
      animate={{
        width: hidden ? "20rem" : "4rem",
        height: "100vh",
      }}
    >
      <motion.nav
        initial={false}
        transition={{ duration: 0.9 }}
        animate={hidden ? "open" : "closed"}
        className="flex flex-col  bg-secondary-gray/20  backdrop-filter backdrop-blur-lg	   shadow-lg h-full"
      >
        <motion.ul variants={variants} className="flex-1 px-3">
          <div className="flex pt-10">
            <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.95 }}>
              <Link href="/" className="flex items-center">
                <Image src={Logo} width={30} height={30} />{" "}
                <motion.span
                  className={`${
                    hidden ? "text-5xl text-white pl-5 " : "scale-0"
                  }`}
                >
                  FinVibe
                </motion.span>
              </Link>
            </motion.div>
            <button type="button" onClick={() => setHidden(!hidden)}>
              <motion.svg
                initial={{ scale: 0 }}
                animate={{ rotate: hidden ? -90 : 90, scale: 1 }}
                transition={{ delay: 0.2 }}
                className="inline w-7 h-7 absolute -right-3 border border-secondary-gray top-8 rounded-full bg-black"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 22 22"
              >
                <defs>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".514"
                      d="m-7 1024.36h34v34h-34z"
                    />
                  </clipPath>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".472"
                      d="m-6 1028.36h32v32h-32z"
                    />
                  </clipPath>
                </defs>
                <path
                  d="m345.44 248.29l-194.29 194.28c-12.359 12.365-32.397 12.365-44.75 0-12.354-12.354-12.354-32.391 0-44.744l171.91-171.91-171.91-171.9c-12.354-12.359-12.354-32.394 0-44.748 12.354-12.359 32.391-12.359 44.75 0l194.29 194.28c6.177 6.18 9.262 14.271 9.262 22.366 0 8.099-3.091 16.196-9.267 22.373"
                  transform="matrix(-.00013-.03541.03541-.00013 3.02 19.02)"
                  fill="#FFFFFF"
                />
              </motion.svg>
            </button>
          </div>
          <div className="pt-5">
            <Links
              link="/dashboard"
              icon={dashboardSvg}
              text="Dashboard"
              hidden={hidden}
              variants={variants}
            />
            <Links
              link="/dashboard/transactions"
              icon={TransactionSvg}
              text="Transactions"
              hidden={hidden}
              variants={variants}
            />
          </div>
        </motion.ul>
        <div className="flex p-3">
          <UserButton />
        </div>
      </motion.nav>
    </motion.aside>
  );
}
