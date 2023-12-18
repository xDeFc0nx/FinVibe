/* eslint-disable import/no-extraneous-dependencies */
import { motion } from "framer-motion";

const Backdrop = ({ children, onClick }) => (
  <motion.div
    onClick={onClick}
    className="bg-black/50 fixed inset-0"
    initial={{ opacity: 0 }}
    animate={{ opacity: 1 }}
    exit={{ opacity: 0 }}
  >
    {children}
  </motion.div>
);

export default Backdrop;
