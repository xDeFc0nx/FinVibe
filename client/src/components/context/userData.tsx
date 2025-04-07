import { createContext, useContext, useState, useEffect } from "react";
import { useWebSocket } from "@/components/WebSocketProvidor";
import type { Description } from "@radix-ui/react-dialog";

export interface UserData {
	ID: string;
	FirstName: string;
	LastName: string;
	Email: string;
	Currency: string;
  Country: string;
}
export interface Account {
	id: string;
	userID: string;
	income: number;
	expense: number;
	balance: number;
	type: string;
}
export interface Transaction {
	ID: string;
	UserID: string;
	AccountID: string;
	Type: string;
	Amount: number;
	Description: string;
	IsRecurring: boolean;
	CreatedAt: string;
}
export interface ChartOverview {
	Day: string;
	Income: number;
	Expense: number;
}
export interface PieOverview {
	Description: string;
	Amount: number;
}
export interface UserDataContextType {
	userData: UserData;
	accounts: Account[];
	activeAccount: Account | null;
	dateRange: string;
	refresh: number;
	chartOverview: ChartOverview[];
	incomePie: PieOverview[];
	expensesPie: PieOverview[];
	transactions: Transaction[];
	setUserData: React.Dispatch<React.SetStateAction<UserData>>;
	setAccounts: React.Dispatch<React.SetStateAction<Account[]>>;
	setDateRange: React.Dispatch<React.SetStateAction<string>>;
	setRefresh: React.Dispatch<React.SetStateAction<number>>;
	setActiveAccount: React.Dispatch<React.SetStateAction<Account | null>>;
	setTransactions: React.Dispatch<React.SetStateAction<Transaction[]>>;
	setChartOverview: React.Dispatch<React.SetStateAction<ChartOverview[]>>;
	setIncomePie: React.Dispatch<React.SetStateAction<PieOverview[]>>;
	setExpensesPie: React.Dispatch<React.SetStateAction<PieOverview[]>>;
}

const UserDataContext = createContext<UserDataContextType | null>(null);

export const useUserData = () => {
	const context = useContext(UserDataContext);
	if (!context) {
		throw new Error("useUserData must be used within a UserDataProvider");
	}
	return context;
};

export const UserDataProvider = ({
	children,
}: { children: React.ReactNode }) => {
	const { socket, isReady } = useWebSocket();
	const [userData, setUserData] = useState<UserData>({
		ID: "",
		FirstName: "",
		LastName: "",
		Email: "",
    Currency: "$",
		Country: "",
	});

	const [accounts, setAccounts] = useState<Account[]>([]);
	const [activeAccount, setActiveAccount] = useState<Account | null>(null);
	const [dateRange, setDateRange] = useState<string>("this_month");
	const [refresh, setRefresh] = useState<number>(0);
	const [transactions, setTransactions] = useState<Transaction[]>([]);
	const [chartOverview, setChartOverview] = useState<ChartOverview[]>([]);
	const [incomePie, setIncomePie] = useState<PieOverview[]>([]);
	const [expensesPie, setExpensesPie] = useState<PieOverview[]>([]);
	useEffect(() => {
		if (socket && isReady) {
			socket.send("getUser");
			socket.send("getAccounts");
			socket.onMessage((msg) => {
				const response = JSON.parse(msg);

				if (response.userData) {
					setUserData(response.userData);
				}
				if (response.accounts) {
					setAccounts(response.accounts);
					if (!activeAccount && response.accounts.length > 0) {
						setActiveAccount(response.accounts[0]);
					}
				}
			});
		}
	}, [socket, isReady]);
	return (
		<UserDataContext.Provider
			value={{
				userData,
				setUserData,
				accounts,
				setAccounts,
				dateRange,
				setDateRange,
				transactions,
				setTransactions,
				activeAccount,
				setActiveAccount,
				chartOverview,
				setChartOverview,
				incomePie,
				setIncomePie,
				expensesPie,
				setExpensesPie,
				refresh,
				setRefresh,
			}}
		>
			{" "}
			{children}
		</UserDataContext.Provider>
	);
};
export const Descriptions = {
	Income: [
		"Salary/Paycheck",
		"Freelance Income",
		"Rental Income",
		"Dividend Payments",
		"Interest Income",
		"Bonus",
		"Sales Revenue",
		"Investment Returns",
		"Government Benefits",
		"Reimbursements",
	],
	Expense: [
		"Groceries",
		"Rent/Mortgage Payment",
		"Utilities",
		"Transportation",
		"Dining Out/Restaurants",
		"Shopping/Retail Purchases",
		"Subscription Services",
		"Insurance",
		"Travel/Vacation Expenses",
		"Medical Bills",
		"Education/Tuition Fees",
		"Loan Repayments",
		"Taxes",
	],
} as const;
