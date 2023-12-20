import React from "react";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";
import Link from "next/link";
import Navbar from "../../../components/ui/Navbar";
import Table from "../../../components/ui/Table";
import Button from "../../../components/ui/Button";

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
          <Link href="/dashboard/new">
            <Button Type="button" color="bg-green-500" Text="New" />
          </Link>
          <Table />
        </div>
      </div>
    </div>
  );
}

export default page;
