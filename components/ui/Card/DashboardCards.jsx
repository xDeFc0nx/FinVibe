// DashboardCards.jsx

"use client";

import React, { useEffect, useState } from "react";
import { CardDashboard } from "./card";
import BalanceSvg from "../../Icon/BalanceSvg.svg";
import IncomeSvg from "../../Icon/IncomeSvg.svg";
import ExpensesSvg from "../../Icon/ExpensesSvg.svg";

function DashboardCards() {
  const [dashboardData, setDashboardData] = useState({});
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await fetch("http://localhost:3000/api/calculations", {
          method: "PUT",
        });

        if (!res.ok) {
          throw new Error("Failed to fetch data");
        }

        const data = await res.json();
        setDashboardData(data);
      } catch (error) {
        console.error(error);
      } finally {
        // Set loading state to false regardless of success or failure
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  if (isLoading) {
    // Display a loading spinner or message
    return <p>Loading...</p>;
  }

  return (
    <div>
      <div className="flex space-x-4 mb-4">
        <CardDashboard
          icon={BalanceSvg}
          header="Balance"
          value={`${dashboardData.sumForBalance.toFixed(2)}$`}
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
