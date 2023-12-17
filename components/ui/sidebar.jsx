/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable jsx-a11y/control-has-associated-label */

"use client";

import { UserButton } from "@clerk/nextjs";
import { useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";

export default function Sidebar() {
  const variants = {
    open: {
      transition: { staggerChildren: 0.07, delayChildren: 0.2 },
    },
    closed: {
      transition: { staggerChildren: 0.05, staggerDirection: -1 },
    },
  };
  const [hidden, setHidden] = useState(true);

  return (
    <motion.aside
      initial={{ width: hidden ? "13rem" : "4rem" }}
      animate={{ width: hidden ? "13rem" : "4rem" }}
      className="h-screen"
    >
      <motion.nav
        initial={false}
        animate={hidden ? "open" : "closed"}
        className="h-full flex flex-col bg-[#24303F]/50  backdrop-filter backdrop-blur-3xl	   shadow-lg"
      >
        <motion.ul variants={variants} className="flex-1 px-3">
          <button
            type="button"
            className="pt-2"
            onClick={() => setHidden(!hidden)}
          >
            <motion.svg
              initial={{ scale: 0 }}
              animate={{ rotate: hidden ? 90 : -90, scale: 1 }}
              className="inline w-5 h-5"
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

          <motion.li
            variants={variants}
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.95 }}
          >
            <Link
              href="/"
              className="flex items-center transition-colors p-2 text-gray-500 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
            >
              <svg
                className="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="currentColor"
                viewBox="0 0 22 21"
              >
                <path d="M16.975 11H10V4.025a1 1 0 0 0-1.066-.998 8.5 8.5 0 1 0 9.039 9.039.999.999 0 0 0-1-1.066h.002Z" />
                <path d="M12.5 0c-.157 0-.311.01-.565.027A1 1 0 0 0 11 1.02V10h8.975a1 1 0 0 0 1-.935c.013-.188.028-.374.028-.565A8.51 8.51 0 0 0 12.5 0Z" />
              </svg>
              <motion.span
                initial={{ opacity: 0 }}
                animate={{ opacity: hidden ? 1 : 0 }}
                transition={{ duration: 0.3 }}
                className={`${hidden ? "ms-3" : "hidden"}`}
              >
                Dashboard
              </motion.span>
            </Link>
          </motion.li>
          <motion.li
            variants={variants}
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.95 }}
          >
            <Link
              href="/"
              className="flex transition-colors items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
            >
              <svg
                className="flex-shrink-0  w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="currentColor"
                viewBox="0 0 18 18"
              >
                <path d="M6.143 0H1.857A1.857 1.857 0 0 0 0 1.857v4.286C0 7.169.831 8 1.857 8h4.286A1.857 1.857 0 0 0 8 6.143V1.857A1.857 1.857 0 0 0 6.143 0Zm10 0h-4.286A1.857 1.857 0 0 0 10 1.857v4.286C10 7.169 10.831 8 11.857 8h4.286A1.857 1.857 0 0 0 18 6.143V1.857A1.857 1.857 0 0 0 16.143 0Zm-10 10H1.857A1.857 1.857 0 0 0 0 11.857v4.286C0 17.169.831 18 1.857 18h4.286A1.857 1.857 0 0 0 8 16.143v-4.286A1.857 1.857 0 0 0 6.143 10Zm10 0h-4.286A1.857 1.857 0 0 0 10 11.857v4.286c0 1.026.831 1.857 1.857 1.857h4.286A1.857 1.857 0 0 0 18 16.143v-4.286A1.857 1.857 0 0 0 16.143 10Z" />
              </svg>
              <motion.span
                initial={{ opacity: 0 }}
                animate={{ opacity: hidden ? 1 : 0 }}
                transition={{ duration: 0.3 }}
                className={`${hidden ? "ms-3" : "hidden"}`}
              >
                kanban
              </motion.span>
            </Link>
          </motion.li>
        </motion.ul>
        <div className="flex p-3">
          <UserButton />
        </div>
      </motion.nav>
    </motion.aside>
  );
}
