/* eslint-disable jsx-a11y/control-has-associated-label */

import React from "react";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";

import Cards from "../components/ui/cards";
import Table from "../components/ui/table";

const page = () => (
  <div>
    <div className="pt-10">
      <ToastContainer />
      <Cards />
      <Table />
    </div>
  </div>
);

export default page;
