import { createContext, useContext, useState, useEffect } from 'react';
import { useWebSocket } from '@/components/WebSocketProvidor';

export interface UserData {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
  Country: string;
}
export interface Account {
  ID: string;
  UserID: string;
  AccountID: string;
  Income: number;
  Expense: number;
  AccountBalance: number;
  Type: string;
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
export interface UserDataContextType {
  userData: UserData;
  accounts: Account[];
  activeAccount: Account | null;
  dateRange: string;
  chartOverview: ChartOverview[];
  transactions: Transaction[];
  setUserData: React.Dispatch<React.SetStateAction<UserData>>;
  setAccounts: React.Dispatch<React.SetStateAction<Account[]>>;
  setDateRange: React.Dispatch<React.SetStateAction<string>>;
  setActiveAccount: React.Dispatch<React.SetStateAction<Account | null>>;
  setTransactions: React.Dispatch<React.SetStateAction<Transaction[]>>;
  setChartOverview: React.Dispatch<React.SetStateAction<ChartOverview[]>>;
}

const UserDataContext = createContext<UserDataContextType | null>(null);

export const useUserData = () => {
  const context = useContext(UserDataContext);
  if (!context) {
    throw new Error('useUserData must be used within a UserDataProvider');
  }
  return context;
};

export const UserDataProvider = ({
  children,
}: { children: React.ReactNode }) => {
  const { socket, isReady } = useWebSocket();
  const [userData, setUserData] = useState<UserData>({
    ID: '',
    FirstName: '',
    LastName: '',
    Email: '',
    Country: '',
  });

  const [accounts, setAccounts] = useState<Account[]>([]);
  const [activeAccount, setActiveAccount] = useState<Account | null>(null);
  const [dateRange, setDateRange] = useState<string>('this_month');
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [chartOverview, setChartOverview] = useState<ChartOverview[]>([]);
  useEffect(() => {
    if (socket && isReady) {
      socket.send('getUser');
      socket.send('getAccounts');
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
      }}
    >
      {' '}
      {children}
    </UserDataContext.Provider>
  );
};
