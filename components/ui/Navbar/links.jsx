import React from "react";
import { motion } from "framer-motion";
import Link from "next/link";
import Image from "next/image";

export default function Links({ link, icon, text, hidden, variants }) {
  return (
    <motion.div>
      <motion.ul>
        <motion.li
          variants={variants}
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.85 }}
        >
          <Link
            href={link}
            className="flex items-center transition-colors p-2 text-gray-500 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-primary-pink group"
          >
            <Image src={icon} className="w-5 h-5" />
            <motion.span
              initial={{ opacity: 0 }}
              animate={{ opacity: hidden ? 1 : 0 }}
              transition={{ duration: 0.3 }}
              className={`${hidden ? "ms-3" : "hidden"}
              text-lg
              `}
            >
              {text}
            </motion.span>
          </Link>
        </motion.li>
      </motion.ul>
    </motion.div>
  );
}
