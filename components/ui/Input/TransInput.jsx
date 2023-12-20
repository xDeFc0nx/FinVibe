"use client";

import { useState } from "react";
import { toast } from "react-toastify";
import { motion } from "framer-motion";
import "react-toastify/dist/ReactToastify.min.css";
import { newTrans } from "../../../actions/addnew";
import Button from "../Button/Button";

import { ALL_CATEGORIES } from "../../../lib/categories";

export default function TransInput() {
  const [inputs, setInputs] = useState({});
  const categories = ALL_CATEGORIES;

  const handleChange = (e) => {
    const { name } = e.target;
    const { value } = e.target;
    setInputs((prevState) => ({ ...prevState, [name]: value }));
  };
  return (
    <div>
      <form
        className="p-4 md:p-5"
        action={async (formData) => {
          setInputs({});
          toast("Added Transaction!", {
            type: "success",
            position: "bottom-right",
          });
          await newTrans(formData);
        }}
      >
        <div className="col-span-2 sm:col-span-1  ">
          <label
            htmlFor="price"
            className="block mb-2 text-sm font-medium text-white "
          >
            Type
          </label>
          <motion.select
            whileFocus={{ scale: 1.1 }}
            required
            value={inputs.type || ""}
            onChange={handleChange}
            placeholder="type"
            name="type"
            className=" border-[2px]  border-secondary-gray text-white bg-secondary-gray/50  backdrop-filter backdrop-blur-lg shadow-lg text-sm rounded-lg block w-full p-2.5 focus:ring-primary-pink focus:border-primary-pink outline-none"
          >
            <option
              selected
              className=" text-white bg-secondary-gray/50  backdrop-filter backdrop-blur-lg shadow-lg text-sm rounded-lg block w-full p-2.5 focus:ring-primary-pink focus:border-primary-pink outline-none"
            >
              Type
            </option>
            {categories &&
              categories.map((category) => (
                <option key={category} value={category}>
                  {category}
                </option>
              ))}
          </motion.select>
        </div>
        <div className="col-span-2 sm:col-span-1">
          <label
            htmlFor="price"
            className="block mb-2 text-sm font-medium text-white "
          >
            Amount
          </label>
          <motion.input
            whileFocus={{ scale: 1.1 }}
            required
            value={inputs.amount || ""}
            onChange={handleChange}
            name="amount"
            type="number"
            className="border-[2px] border-secondary-gray text-white bg-secondary-gray/50  backdrop-filter backdrop-blur-lg shadow-lg text-sm rounded-lg block w-full p-2.5 focus:ring-primary-pink focus:border-primary-pink outline-none"
            placeholder="$2999"
          />
        </div>
        <div className="col-span-2">
          <label
            htmlFor="description"
            className="block mb-2 text-sm font-medium text-white "
          >
            Description
          </label>
          <motion.textarea
            whileFocus={{ scale: 1.1 }}
            required
            value={inputs.description || ""}
            onChange={handleChange}
            name="description"
            rows="4"
            className="border-[2px] border-secondary-gray text-white bg-secondary-gray/50  backdrop-filter backdrop-blur-lg shadow-lg text-sm rounded-lg block w-full p-2.5 focus:ring-primary-pink focus:border-primary-pink outline-none"
            placeholder="Write description here"
          />
        </div>
        <div className="pt-7">
          <Button Text="Add New" color="bg-green-500" type="submit" />
        </div>
      </form>
    </div>
  );
}
