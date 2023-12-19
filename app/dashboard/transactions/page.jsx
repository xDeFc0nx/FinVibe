import React from "react";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";
import Navbar from "../../../components/ui/Navbar";
import Table from "../../../components/ui/Table";
import NewBtn from "../../../components/ui/Table/newTransaction";

function page() {
  return (
    <div className="flex h-screen">
      <ToastContainer />
      {/* Navbar */}
      <div className="flex-none">
        <Navbar />
      </div>

      {/* Main Content */}
      <div className="flex flex-col flex-1 p-20">
        <div className="grid grid-cols-1 gap-5 pt-10">
          <NewBtn />
          <Table />
        </div>
      </div>
    </div>
  );
}

export default page;
