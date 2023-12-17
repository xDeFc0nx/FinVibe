import { NextResponse } from "next/server";
import { revalidatePath } from "next/cache";
import prisma from "../../../../prisma/client";

export const revalidate = true;

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
    console.error("Error fetching transaction:", err);
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

    revalidatePath(`/api/transactions/${id}`); // Revalidate cache for this specific path

    return NextResponse.json(updateTransaction);
  } catch (err) {
    console.error("Error updating transaction:", err);
    return NextResponse.json(
      { message: "Error updating transaction", err },
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

    revalidatePath(`/api/transactions/${id}`); // Revalidate cache for this specific path

    return NextResponse.json(transaction);
  } catch (err) {
    console.error("Error deleting transaction:", err);
    return NextResponse.json(
      { message: "Error deleting transaction", err },
      { status: 500 }
    );
  }
};
