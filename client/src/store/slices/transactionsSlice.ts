import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { Transaction, LoadingStatus } from '@/types';


export interface TransactionsState {
  list: Transaction[];
  dateRange: string;
  status: LoadingStatus;
  error: string | null;
}

const initialState: TransactionsState = {
  list: [],
  dateRange: 'this_month',
  status: 'idle',
  error: null,
};

const transactionsSlice = createSlice({
  name: 'transactions',
  initialState,
  reducers: {
    transactionsLoading(state) {
      state.status = 'loading';
      state.error = null;
    },
    transactionsReceived(state, action: PayloadAction<Transaction[]>) {
      state.list = action.payload;
      state.status = 'succeeded';
    },
    addTransaction(state, action: PayloadAction<Transaction>) {
      const exists = state.list.some(tx => tx.ID === action.payload.ID);
      if (!exists) {
        state.list.unshift(action.payload);
      }
    },
    removeTransaction(state, action: PayloadAction<string>) {
      state.list = state.list.filter(tx => tx.ID !== action.payload);
      state.status = 'succeeded';
    },
    setDateRange(state, action: PayloadAction<string>) {
      state.dateRange = action.payload;
      state.status = 'idle';
      state.list = [];
    },
    transactionsError(state, action: PayloadAction<string>) {
      state.status = 'failed';
      state.error = action.payload;
      state.list = [];
    },
    clearTransactions(state) {
      Object.assign(state, initialState);
    }
  },
});

export const {
  transactionsLoading,
  transactionsReceived,
  addTransaction,
  removeTransaction,
  setDateRange,
  transactionsError,
  clearTransactions
} = transactionsSlice.actions;
export default transactionsSlice.reducer;
