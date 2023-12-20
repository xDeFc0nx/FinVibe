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
      transactions.find((item) => Income.includes(item.type))?._sum.amount || 0;

    const sumForExpense =
      transactions.find((item) => Expenses.includes(item.type))?._sum.amount ||
      0;
    const sumForSpecial =
      transactions.find((item) => Special.includes(item.type))?._sum.amount ||
      0;

    const sumForBalance = sumForIncome + sumForExpense + sumForSpecial || 0;
    return {
      sumForIncome,
      sumForExpense,
      sumForBalance,
    };
  } catch (err) {
    console.error(err);
  }
};
