// url: 'http://localhost:3000/api/transactions/12312

import { NextResponse } from "next/server";
import prisma from "../../../../prisma/client";

export const GET = async (request, { params }) => {
  try {
    const { id } = params;
    const transaction = await prisma.transaction.findUnique({
      where: {
        id,
      },
    });
    if (!transaction) {
      return NextResponse.json(
        { message: "Transaction not found" },
        { status: 404 }
      );
    }

    return NextResponse.json(transaction);
  } catch (err) {
    console.error("Error creating transaction:", err);
    return NextResponse.json({ message: "GET error", err }, { status: 500 });
  }
};

export const PATCH = async (request, { params }) => {
  try {
    const body = await request.json();
    const { type, amount, description } = body;
    const { id } = params;

    const updateTransaction = await prisma.transaction.update({
      where: {
        id,
      },
      data: {
        type,
        amount,
        description,
      },
    });
    if (!updateTransaction) {
      return NextResponse.json(
        { message: "Transaction not found" },
        { status: 404 }
      );
    }
    return NextResponse.json(updateTransaction);
  } catch (err) {
    return NextResponse.json(
      { message: "Error Updating transaction", err },
      { status: 500 }
    );
  }
};

export const DELETE = async (request, { params }) => {
  try {
    const { id } = params;
    const transaction = await prisma.transaction.delete({
      where: {
        id,
      },
    });
    if (!transaction) {
      return NextResponse.json(
        { message: "Transaction not found" },
        { status: 404 }
      );
    }

    return NextResponse.json(transaction);
  } catch (err) {
    console.error("Error Deleting transaction:", err);
    return NextResponse.json({ message: "GET error", err }, { status: 500 });
  }
};
