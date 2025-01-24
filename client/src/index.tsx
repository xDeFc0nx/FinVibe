import './App.css';
import App from '@/App';
import Index from '@/pages/(secure)/index';
import Settings from '@/pages/(secure)/settings';
import Transactions from '@/pages/(secure)/transactions';
import { Layout } from '@/pages/(secure)/layout';
import Login from '@/pages/login';
import Register from '@/pages/register';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router';
import { ToastContainer } from 'react-toastify';
const root = document.getElementById('root');

// @ts-ignore
ReactDOM.createRoot(root).render(
  <BrowserRouter>
    <ToastContainer />
    <Routes>
      <Route path="/" element={<App />} />
      <Route path="login" element={<Login />} />
      <Route path="register" element={<Register />} />
      <Route path="/app" element={<Layout />}>
        <Route path="dashboard" element={<Index />} />
        <Route path="settings" element={<Settings />} />
        <Route path="transactions" element={<Transactions />} />
      </Route>
    </Routes>
  </BrowserRouter>,
);
