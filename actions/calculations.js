/* eslint-disable no-underscore-dangle */
// calculations.js

"use server";

import prisma from "../prisma/client";

import {
  SPECIAL_CATEGORIES,
  EXPENSES_CATEGORIES,
  INCOMES_CATEGORIES,
  ALL_CATEGORIES,
} from "../lib/categories";

export const Calculations = async () => {
  try {
    const types = ALL_CATEGORIES;
    const Income = INCOMES_CATEGORIES;
    const Expenses = EXPENSES_CATEGORIES;
    const Special = SPECIAL_CATEGORIES;
    const Investing = "Investing";
    const Savings = "Savings";

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

    const sumForIncome = transactions
      .filter((item) => Income.includes(item.type))
      .reduce((acc, cur) => acc + cur._sum.amount, 0);

    const sumForExpense = transactions
      .filter((item) => Expenses.includes(item.type))
      .reduce((acc, cur) => acc + cur._sum.amount, 0);

    const sumForSpecial = transactions
      .filter((item) => Special.includes(item.type))
      .reduce((acc, cur) => acc + cur._sum.amount, 0);

    const sumForSavings = transactions
      .filter((item) => Savings.includes(item.type))
      .reduce((acc, cur) => acc + cur._sum.amount, 0);

    const sumForInvesting = transactions
      .filter((item) => Investing.includes(item.type))
      .reduce((acc, cur) => acc + cur._sum.amount, 0);

    const sumForBalance = sumForIncome + sumForExpense + sumForSpecial;

    return {
      sumForIncome,
      sumForExpense,
      sumForSavings,
      sumForInvesting,
      sumForBalance,
    };
  } catch (err) {
    console.error(err);
    return { error: "Error fetching calculations" };
  }
};
