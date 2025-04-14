import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { ChartOverview, PieOverview, LoadingStatus } from '@/types';

export interface OverviewState {
  chartData: ChartOverview[];
  incomePieData: PieOverview[];
  expensePieData: PieOverview[];
  status: LoadingStatus;
  error: string | null;
}

const initialState: OverviewState = {
  chartData: [],
  incomePieData: [],
  expensePieData: [],
  status: 'idle',
  error: null,
};

const overviewSlice = createSlice({
  name: 'overview',
  initialState,
  reducers: {
    overviewLoading(state) {
      state.status = 'loading';
      state.error = null;
    },
    overviewReceived(state, action: PayloadAction<{ chart: ChartOverview[], incomePie: PieOverview[], expensePie: PieOverview[] }>) {
      state.chartData = action.payload.chart;
      state.incomePieData = action.payload.incomePie;
      state.expensePieData = action.payload.expensePie;
      state.status = 'succeeded';
    },
    overviewError(state, action: PayloadAction<string>) {
      state.status = 'failed';
      state.error = action.payload;
      state.chartData = [];
      state.incomePieData = [];
      state.expensePieData = [];
    },
    clearOverview(state) {
      Object.assign(state, initialState);
    }
  },
});

export const {
  overviewLoading,
  overviewReceived,
  overviewError,
  clearOverview
} = overviewSlice.actions;
export default overviewSlice.reducer;
