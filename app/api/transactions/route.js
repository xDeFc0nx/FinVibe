// url: 'http://localhost:3000/api/transactions

import { NextResponse } from "next/server";
import prisma from "../../../prisma/client";

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
    return NextResponse.json(newTransaction);
  } catch (err) {
    return NextResponse.json(
      { message: "Error creating transaction", err },
      { status: 500 }
    );
  }
};

export const GET = async () => {
  try {
    const transactions = await prisma.transaction.findMany({
      select: {
        id: true,
        type: true,
        amount: true,
        description: true,
      },
    });
    return NextResponse.json(transactions);
  } catch (err) {
    return NextResponse.json(
      { message: "Error Feching transaction", err },
      { status: 500 }
    );
  }
};
