import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { Account, LoadingStatus } from '../../types';

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
  accountsLoading,
  accountsReceived,
  setActiveAccount,
  accountsError,
  clearAccounts
} = accountsSlice.actions;
export default accountsSlice.reducer;
