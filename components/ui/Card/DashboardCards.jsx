// DashboardCards.jsx

"use client";

import React, { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { CardDashboard } from "./card";
import BalanceSvg from "../../Icon/BalanceSvg.svg";
import IncomeSvg from "../../Icon/IncomeSvg.svg";
import ExpensesSvg from "../../Icon/ExpensesSvg.svg";
import { Calculations } from "../../../actions/calculations";

function DashboardCards() {
  const [dashboardData, setDashboardData] = useState({});
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await Calculations();
        setDashboardData(data);
        setIsLoading(false);
      } catch (error) {
        console.error(error);
        toast.error(error.message);
      }
    };

    fetchData();
  }, []);

  if (isLoading) {
    return <p>Loading...</p>;
  }

  return (
    <div>
      <div className="flex space-x-4 mb-4">
        <CardDashboard
          icon={BalanceSvg}
          header="Balance"
          value={`+${dashboardData.sumForBalance.toFixed(2)}$`}
          color="text-white"
        />
        <CardDashboard
          icon={IncomeSvg}
          header="Income"
          value={`+${dashboardData.sumForIncome.toFixed(2)}$`}
          color="text-green-500"
        />
        <CardDashboard
          icon={ExpensesSvg}
          header="Expenses"
          value={`-${dashboardData.sumForExpense.toFixed(2)}$`}
          color="text-red-500"
        />
      </div>
    </div>
  );
}

export default DashboardCards;
