
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

export type LoadingStatus = 'idle' | 'loading' | 'succeeded' | 'failed';
