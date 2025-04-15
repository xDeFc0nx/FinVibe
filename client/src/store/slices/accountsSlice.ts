import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { Account, LoadingStatus } from '@/types';
interface AccountUpdatePayload {
  id: string;
  details: {
    Balance: number;
    Income: number;
    Expense: number;
  };
}
export interface AccountsState {
  list: Account[];
  activeAccountId: string | null;
  status: LoadingStatus;
  error: string | null;
}

const initialState: AccountsState = {
  list: [],
  activeAccountId: null,
  status: 'idle',
  error: null,
};

const accountsSlice = createSlice({
  name: 'accounts',
  initialState,
  reducers: {
    addAccount(state, action: PayloadAction<Account>) {
      const exists = state.list.some(acc => acc.ID === action.payload.ID);
      if (!exists) {
        state.list.push(action.payload);
      }
      state.status = 'succeeded';
    },
    accountsLoading(state) {
      state.status = 'loading';
      state.error = null;
    },
    accountsReceived(state, action: PayloadAction<Account[]>) {
      state.list = action.payload;
      state.status = 'succeeded';
      if (!state.activeAccountId && action.payload.length > 0) {
        state.activeAccountId = action.payload[0].ID;
      }
      if (state.activeAccountId && !action.payload.find(acc => acc.ID === state.activeAccountId)) {
        state.activeAccountId = action.payload.length > 0 ? action.payload[0].ID : null;
      }
    },

    setActiveAccount(state, action: PayloadAction<string | null>) {
      state.activeAccountId = action.payload;
    },
    updateAccountDetails(state, action: PayloadAction<AccountUpdatePayload>) {
      const { id, details } = action.payload;
      const accountIndex = state.list.findIndex(acc => acc.ID === id);

      if (accountIndex !== -1) {
        state.list[accountIndex] = {
          ...state.list[accountIndex],
          ...details,
        };
      } else {
        console.warn(`Account with ID ${id} not found in state to update.`);
      }
    },
    accountsError(state, action: PayloadAction<string>) {
      state.status = 'failed';
      state.error = action.payload;
      state.list = [];
      state.activeAccountId = null;
    },
    clearAccounts(state) {
      Object.assign(state, initialState);
    }
  },
});

export const {
  addAccount,
  accountsLoading,
  accountsReceived,
  updateAccountDetails,
  setActiveAccount,
  accountsError,
  clearAccounts
} = accountsSlice.actions;
export default accountsSlice.reducer;
