"use server";

import { revalidatePath } from "next/cache";
import prisma from "../prisma/client";

export const getId = async ({ params }) => {
  try {
    const { id } = params;
    const transaction = await prisma.transaction.findUnique({
      where: {
        id,
      },
    });
    if (!transaction) {
      console.error("Couldn't find transaction");
    }
    revalidatePath("/dashboard/transactions");

    return { status: 200, body: transaction };
  } catch (err) {
    console.error("Error fetching transaction:", err);
  }
};

export const DeleteData = async ({ params }) => {
  try {
    const { id } = params;

    if (!id) {
      return { message: "Missing transaction ID", status: 400 };
    }

    const transaction = await prisma.transaction.delete({
      where: {
        id,
      },
    });

    if (!transaction) {
      return { message: "Transaction not found", status: 404 };
    }

    revalidatePath("/dashboard/transactions");

    return { message: "Transaction deleted successfully" };
  } catch (err) {
    console.error("Error deleting transaction:", err);
    return { message: "Error deleting transaction", err, status: 500 };
  }
};
