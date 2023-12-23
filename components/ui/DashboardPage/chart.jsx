/* eslint-disable import/no-extraneous-dependencies */

"use client";

import React, { useEffect, useState } from "react";

import { toast } from "react-toastify";
import { LineChart, Line, XAxis, YAxis, ResponsiveContainer } from "recharts";
import { CardChart } from "../Card/card";
import { Calculations } from "../../../actions/calculations";

export function Chart() {
  const [chartData, setChartData] = useState({});
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await Calculations();
        setChartData(data);
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

  const data = Object.values(chartData).map((value, index) => ({
    uv: value,
    day: index + 1,
  }));

  return (
    <>
      <CardChart header="Balance" styles="h-40 w-full">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={data}>
            <XAxis dataKey="day" />
            <YAxis />

            <Line type="monotone" dataKey="uv" stroke="#8884d8" />
          </LineChart>
        </ResponsiveContainer>
      </CardChart>
      <div className="flex gap-2">
        <CardChart header="Balance" styles="text-white" />
        <CardChart header="Balance" styles="text-white" />
        <CardChart header="Balance" styles="text-white" />
      </div>
    </>
  );
}
