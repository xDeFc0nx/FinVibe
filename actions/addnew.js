"use server";

import { revalidatePath } from "next/cache";

import prisma from "../prisma/client";

export const newTrans = async (formData) => {
  const type = formData.get("type");
  const amount = parseInt(formData.get("amount"), 10);
  const description = formData.get("description");

  await prisma.transaction.create({
    data: {
      type,
      amount,
      description,
    },
  });

  revalidatePath("/dashboard/transactions");
};
