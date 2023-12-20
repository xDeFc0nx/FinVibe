"use server";

import prisma from "../prisma/client";

export const FetchData = async () => {
  try {
    const transactions = await prisma.transaction.findMany({
      select: {
        id: true,
        type: true,
        amount: true,
        description: true,
        DateCreated: true,
      },
    });
    return transactions;
  } catch (err) {
    console.error(err);
    throw err; // Rethrow the error so that it can be caught by the component
  }
};
