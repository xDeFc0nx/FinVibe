import React from "react";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";
import Navbar from "../../components/ui/Navbar";
import Card from "../../components/ui/DashboardPage/DashboardCards";
import Topnav from "../../components/ui/DashboardPage/topnav";
import { Chart } from "../../components/ui/DashboardPage/chart";

const IndexPage = () => (
  <div className="flex h-screen">
    <ToastContainer />
    {/* Navbar */}
    <div className="flex-none">
      <Navbar />
    </div>

    {/* Main Content */}
    <div className="flex flex-col flex-1 p-20">
      <Topnav />
      <Card />

      <div className="grid grid-cols-1  gap-5 pt-10">
        <Chart />
      </div>
    </div>
  </div>
);

export default IndexPage;
