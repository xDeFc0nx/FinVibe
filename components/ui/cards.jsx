/* eslint-disable import/no-extraneous-dependencies */
/* eslint no-underscore-dangle: 0 */

import React from "react";
import { Card, CardHeader, CardBody } from "@nextui-org/card";
import Image from "next/image";
import Wallet3 from "../Icon/wallet-3.svg";
import Wallet from "../Icon/wallet.svg";
import cardsvg from "../Icon/card.svg";
import dollar from "../Icon/dollar-circle.svg";
import prisma from "../../prisma/client";

export default async function Cards() {
  const types = ["Balance", "Income", "Expense", "Savings"];

  const transactions = await prisma.transaction.groupBy({
    by: ["type"],
    _sum: {
      amount: true,
    },
    where: {
      type: {
        in: types,
      },
    },
  });

  const sumForIncome =
    transactions.find((item) => item.type === "Income")?._sum.amount || 0;
  const sumForExpense =
    transactions.find((item) => item.type === "Expense")?._sum.amount || 0;
  const sumForSavings =
    transactions.find((item) => item.type === "Savings")?._sum.amount || 0;
  const sumForBalance = sumForIncome + sumForExpense || 0;

  return (
    <div className="grid grid-cols-4 space-x-10 backdrop-blur-3xl    ">
      <Card className="w-[18rem]  text-white  bg-[#24303F]/50  backdrop-filter backdrop-blur-3xl	   shadow-lg">
        <CardHeader className="w-full justify-center">
          {/* Image on the left */}
          <Image
            alt="wallet"
            height={30}
            width={30}
            radius="sm"
            src={Wallet3}
          />

          {/* Text in the center */}
          <div className="flex flex-col justify-center">
            <p className="text-md  text-gray-400 text-[1.2rem]">Balance</p>
          </div>
        </CardHeader>

        <CardBody>
          <h1
            className={`${
              sumForBalance >= 0 ? "text-green-600" : "text-red-600"
            } text-[2.5rem] text-center`}
          >
            ${sumForBalance}
            <span
              className={`${
                sumForBalance >= 0 ? "bg-green-600" : "bg-red-600"
              } text-sm  text-white  rounded  w-10 flex text-center `}
            >
              <svg
                className={`${
                  sumForBalance >= 0 ? "flex" : "rotate-180"
                } inline w-5 h-5 `}
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 22 22"
              >
                <defs>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".514"
                      d="m-7 1024.36h34v34h-34z"
                    />
                  </clipPath>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".472"
                      d="m-6 1028.36h32v32h-32z"
                    />
                  </clipPath>
                </defs>
                <path
                  d="m345.44 248.29l-194.29 194.28c-12.359 12.365-32.397 12.365-44.75 0-12.354-12.354-12.354-32.391 0-44.744l171.91-171.91-171.91-171.9c-12.354-12.359-12.354-32.394 0-44.748 12.354-12.359 32.391-12.359 44.75 0l194.29 194.28c6.177 6.18 9.262 14.271 9.262 22.366 0 8.099-3.091 16.196-9.267 22.373"
                  transform="matrix(-.00013-.03541.03541-.00013 3.02 19.02)"
                  fill="#FFFFFF"
                />
              </svg>
              1.5%
            </span>
          </h1>
        </CardBody>
      </Card>
      <Card className="w-[18rem]  text-white bg-[#24303F]/50  backdrop-filter backdrop-blur-3xl	   shadow-lg ">
        <CardHeader className="w-full justify-center">
          {/* Image on the left */}
          <Image alt="wallet" height={30} width={30} radius="sm" src={Wallet} />

          {/* Text in the center */}
          <div className="flex flex-col justify-center">
            <p className="text-md  text-gray-400 text-[1.2rem]">Income</p>
          </div>
        </CardHeader>

        <CardBody>
          <h1
            className={`${
              sumForIncome >= 0 ? "text-green-600" : "text-red-600"
            } text-[2.5rem] text-center`}
          >
            ${sumForIncome}
            <span
              className={`${
                sumForIncome >= 0 ? "bg-green-600" : "bg-red-600"
              } text-sm  text-white  rounded  w-10 flex text-center `}
            >
              <svg
                className={`${
                  sumForIncome >= 0 ? "flex" : "rotate-180"
                } inline w-5 h-5 `}
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 22 22"
              >
                <defs>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".514"
                      d="m-7 1024.36h34v34h-34z"
                    />
                  </clipPath>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".472"
                      d="m-6 1028.36h32v32h-32z"
                    />
                  </clipPath>
                </defs>
                <path
                  d="m345.44 248.29l-194.29 194.28c-12.359 12.365-32.397 12.365-44.75 0-12.354-12.354-12.354-32.391 0-44.744l171.91-171.91-171.91-171.9c-12.354-12.359-12.354-32.394 0-44.748 12.354-12.359 32.391-12.359 44.75 0l194.29 194.28c6.177 6.18 9.262 14.271 9.262 22.366 0 8.099-3.091 16.196-9.267 22.373"
                  transform="matrix(-.00013-.03541.03541-.00013 3.02 19.02)"
                  fill="#FFFFFF"
                />
              </svg>
              1.5%
            </span>
          </h1>
        </CardBody>
      </Card>
      <Card className="w-[18rem] ] text-white  bg-[#24303F]/50  backdrop-filter backdrop-blur-3xl	   shadow-lg">
        <CardHeader className="w-full justify-center">
          {/* Image on the left */}
          <Image alt="card" height={30} width={30} radius="sm" src={cardsvg} />

          {/* Text in the center */}
          <div className="flex flex-col justify-center">
            <p className="text-md  text-gray-400 text-[1.2rem]">Expenses</p>
          </div>
        </CardHeader>

        <CardBody>
          <h1 className="text-[2.5rem] text-center text-red-600">
            ${sumForExpense}
            <span
              className={`${
                sumForExpense >= 0 ? "bg-green-600" : "bg-red-600"
              } text-sm  text-white  rounded  w-10 flex text-center `}
            >
              <svg
                className={`${
                  sumForExpense >= 0 ? "flex" : "rotate-180"
                } inline w-5 h-5 `}
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 22 22"
              >
                <defs>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".514"
                      d="m-7 1024.36h34v34h-34z"
                    />
                  </clipPath>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".472"
                      d="m-6 1028.36h32v32h-32z"
                    />
                  </clipPath>
                </defs>
                <path
                  d="m345.44 248.29l-194.29 194.28c-12.359 12.365-32.397 12.365-44.75 0-12.354-12.354-12.354-32.391 0-44.744l171.91-171.91-171.91-171.9c-12.354-12.359-12.354-32.394 0-44.748 12.354-12.359 32.391-12.359 44.75 0l194.29 194.28c6.177 6.18 9.262 14.271 9.262 22.366 0 8.099-3.091 16.196-9.267 22.373"
                  transform="matrix(-.00013-.03541.03541-.00013 3.02 19.02)"
                  fill="#FFFFFF"
                />
              </svg>
              1.5%
            </span>
          </h1>
        </CardBody>
      </Card>
      <Card className="w-[18rem]  text-white bg-[#24303F]/50  backdrop-filter backdrop-blur-3xl	   shadow-lg ">
        <CardHeader className="w-full justify-center">
          {/* Image on the left */}
          <Image alt="dollar" height={30} width={30} radius="sm" src={dollar} />

          {/* Text in the center */}
          <div className="flex flex-col justify-center">
            <p className="text-md  text-gray-400 text-[1.2rem]">Savings</p>
          </div>
        </CardHeader>

        <CardBody>
          <h1
            className={`${
              sumForSavings >= 0 ? "text-green-600" : "text-red-600"
            } text-[2.5rem] text-center`}
          >
            ${sumForSavings}
            <span
              className={`${
                sumForSavings >= 0 ? "bg-green-600" : "bg-red-600"
              } text-sm  text-white  rounded  w-10 flex text-center `}
            >
              <svg
                className={`${
                  sumForSavings >= 0 ? "flex" : "rotate-180"
                } inline w-5 h-5 `}
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 22 22"
              >
                <defs>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".514"
                      d="m-7 1024.36h34v34h-34z"
                    />
                  </clipPath>
                  <clipPath>
                    <path
                      fill="#FFFFFF"
                      fillOpacity=".472"
                      d="m-6 1028.36h32v32h-32z"
                    />
                  </clipPath>
                </defs>
                <path
                  d="m345.44 248.29l-194.29 194.28c-12.359 12.365-32.397 12.365-44.75 0-12.354-12.354-12.354-32.391 0-44.744l171.91-171.91-171.91-171.9c-12.354-12.359-12.354-32.394 0-44.748 12.354-12.359 32.391-12.359 44.75 0l194.29 194.28c6.177 6.18 9.262 14.271 9.262 22.366 0 8.099-3.091 16.196-9.267 22.373"
                  transform="matrix(-.00013-.03541.03541-.00013 3.02 19.02)"
                  fill="#FFFFFF"
                />
              </svg>
              1.5%
            </span>
          </h1>
        </CardBody>
      </Card>
    </div>
  );
}
