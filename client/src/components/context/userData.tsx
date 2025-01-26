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
  DateRange: string;
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
export interface UserDataContextType {
  userData: UserData;
  accounts: Account[];
  activeAccount: Account | null;
  transactions: Transaction[];
  setUserData: React.Dispatch<React.SetStateAction<UserData>>;
  setAccounts: React.Dispatch<React.SetStateAction<Account[]>>;
  setActiveAccount: React.Dispatch<React.SetStateAction<Account | null>>;
  setTransactions: React.Dispatch<React.SetStateAction<Transaction[]>>;
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

  const [transactions, setTransactions] = useState<Transaction[]>([]);

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
        transactions,
        setTransactions,
        activeAccount,
        setActiveAccount,
      }}
    >
      {' '}
      {children}
    </UserDataContext.Provider>
  );
};
