import React from "react";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";
import Navbar from "../../components/ui/Navbar";
import Card from "../../components/ui/Card/DashboardCards";

const IndexPage = () => (
  <div className="flex h-screen">
    <ToastContainer />
    {/* Navbar */}
    <div className="flex-none">
      <Navbar />
    </div>

    {/* Main Content */}
    <div className="flex flex-col flex-1 p-20">
      <Card />

      <div className="grid grid-cols-1 gap-5 pt-10">
        {/* <div className="w-full h-full  bg-secondary-gray/50  backdrop-filter backdrop-blur-lg	    shadow-lg">
          test
        </div> */}
      </div>
    </div>
  </div>
);

export default IndexPage;
