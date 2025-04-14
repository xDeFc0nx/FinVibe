import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { UserData, LoadingStatus } from '@/types';

export interface UserState {
  data: UserData | null;
  status: LoadingStatus;
  error: string | null;
}

const initialState: UserState = {
  data: null,
  status: 'idle',
  error: null,
};

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    userLoading(state: any) {
      state.status = 'loading';
      state.error = null;
    },
    userReceived(state: any, action: PayloadAction<UserData>) {
      state.data = action.payload;
      state.status = 'succeeded';
    },
    userError(state: any, action: PayloadAction<string>) {
      state.status = 'failed';
      state.error = action.payload;
      state.data = null;
    },
    clearUser(state: any) {
      Object.assign(state, initialState); // Reset state
    }
  },
});

export const {
  userLoading,
  userReceived,
  userError,
  clearUser
} = userSlice.actions;
export default userSlice.reducer;
