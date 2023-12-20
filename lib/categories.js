export const ExpenseCategories = [
  "Home",
  "Groceries",
  "Transportation",
  "Subscriptions",
  "Trips",
  "Hobbies",
  "Health",
  "Bar_Cafe_Restaurant",
  "Clothes & Shoes",
  "Internet",
  "Others",
  "Online Shopping",
  "Donations & Gifts",
  "Rent",
];

export const IncomeCategories = ["Salary", "Refunds / Reimbursements"];

export const SpecialCategories = ["Savings", "Investing", "Uncategorized"];

export const EXPENSES_CATEGORIES = Object.values(ExpenseCategories).sort();
export const INCOMES_CATEGORIES = Object.values(IncomeCategories).sort();
export const SPECIAL_CATEGORIES = Object.values(SpecialCategories).sort();

export const ALL_CATEGORIES = [
  ...EXPENSES_CATEGORIES,
  ...INCOMES_CATEGORIES,
  ...SPECIAL_CATEGORIES,
];
