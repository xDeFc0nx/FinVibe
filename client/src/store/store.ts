import { configureStore } from '@reduxjs/toolkit';
import userReducer from './slices/userSlice';
import accountsReducer from './slices/accountsSlice';
import transactionsReducer from './slices/transactionsSlice';
import overviewReducer from './slices/chartsSlice';


export const store = configureStore({
  reducer: {
    user: userReducer,
    accounts: accountsReducer,
    transactions: transactionsReducer,
    overview: overviewReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
