import React from "react";
import RealTimeSvg from "@/assets/RealTime.svg";
import PersonalizedSvg from "@/assets/Personalized.svg";
import SecureSvg from "@/assets/Secure.svg";
import ExpenseSvg from "@/assets/Expense.svg";
import Card from "./cards"
const CardsHero = () => (
  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 md:gap-8">
    <Card
      icon={RealTimeSvg}
      header="Real-time Financial Insights"
      value="Display real-time updates on financial data, providing users with the latest information on their accounts, transactions, and investment performance."
    />
    <Card
      icon={PersonalizedSvg}
      header="Personalized Budgeting"
      value="Ensure the security of financial transactions by offering a feature that monitors and alerts users about potentially suspicious or unauthorized activities."
    />
    <Card
      icon={SecureSvg}
      header="Secure Transaction Monitoring"
      value="Provide users with detailed analytics on their investments, including performance charts, historical data, and recommendations for optimizing their investment portfolio."
    />
    <Card
      icon={ExpenseSvg}
      header="Expense Categorization"
      value="Automatically categorize expenses to help users understand their spending habits better. This feature can simplify budgeting and enable users to identify areas where they can save money."
    />
  </div>
);

export default CardsHero;
