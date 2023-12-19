import { NextResponse } from "next/server";
import { revalidatePath } from "next/cache";
import prisma from "../../../prisma/client";

export const revalidate = true;

export const POST = async (request) => {
  try {
    const body = await request.json();
    const { type, amount, description } = body;

    const newTransaction = await prisma.transaction.create({
      data: {
        type,
        amount,
        description,
      },
    });
    revalidatePath("/api/transactions");

    // Display success toast
    return NextResponse.json(newTransaction, {
      revalidate: 3,
    });
  } catch (err) {
    // Display error toast
    console.error(err);
    return NextResponse.json(
      { message: "Error creating transaction", err },
      { status: 500 }
    );
  }
};

export const PUT = async () => {
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

    // Display success toast
    return NextResponse.json(transactions, {
      revalidate: 3, // Revalidate every 3 seconds
    });
  } catch (err) {
    // Display error toast
    console.error(err);
    return NextResponse.json(
      { message: "Error fetching transactions", err },
      { status: 500 }
    );
  }
};
