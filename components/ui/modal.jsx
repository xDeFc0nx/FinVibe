/* eslint-disable jsx-a11y/control-has-associated-label */
/* eslint-disable react/jsx-no-useless-fragment */
/* eslint-disable import/no-extraneous-dependencies */
import React, { useEffect } from "react";
import { motion, useAnimation } from "framer-motion";
import BackDrop from "./backdrop";

const Modal = ({ children, modalOpen, text, handleClose }) => {
  const zoomIn = {
    type: "spring",
    damping: 25,
    stiffness: 120,
  };

  const modalAnimation = {
    hidden: { opacity: 0, scale: 0.5 },
    visible: { opacity: 1, scale: 1, transition: zoomIn },
    exit: { opacity: 0, scale: 0.5, transition: { duration: 0.2 } },
  };

  const controls = useAnimation();

  useEffect(() => {
    if (modalOpen) {
      controls.start("visible");
    } else {
      controls.start("exit");
    }
  }, [modalOpen, controls]);

  return (
    <>
      {modalOpen && (
        <BackDrop onClick={handleClose}>
          <motion.div onClick={(e) => e.stopPropagation()}>
            <motion.div
              initial="hidden"
              animate={controls}
              exit="exit"
              variants={modalAnimation}
              className="flex justify-center items-center h-full rounded-lg"
            >
              <motion.div className="items-end bg-[#24303F]/50  backdrop-filter backdrop-blur-3xl	   shadow-lg rounded-lg  w-1/2 h-1/2 p-5">
                <motion.div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-600">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                    {text}
                  </h3>

                  <button
                    type="button"
                    className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white"
                    onClick={handleClose}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M6 18L18 6M6 6l12 12"
                      />
                    </svg>
                  </button>
                </motion.div>
                {children}
              </motion.div>
            </motion.div>
          </motion.div>
        </BackDrop>
      )}
    </>
  );
};

export default Modal;
