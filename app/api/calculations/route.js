/* eslint-disable no-underscore-dangle */
// calculations.js

import { NextResponse } from "next/server";
import prisma from "../../../prisma/client";

export const PUT = async () => {
  try {
    const types = ["Income", "Expense"];

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
    const sumForBalance = sumForIncome - sumForExpense || 0;

    return NextResponse.json({
      sumForIncome,
      sumForExpense,
      sumForBalance,
      revalidate: 3,
    });
  } catch (err) {
    console.error(err);
    return NextResponse.json(
      { message: "Error fetching transactions", err },
      { status: 500 }
    );
  }
};
