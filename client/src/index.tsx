import "./App.css";
import App from "@/App";
import Index from "@/pages/(secure)/index";
import { Layout } from "@/pages/(secure)/layout";
import Login from "@/pages/login";
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
const root = document.getElementById("root");

// @ts-ignore
ReactDOM.createRoot(root).render(
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<App />} />
      <Route path="/login" element={<Login />} />
      <Route path="/app" element={<Layout />}>
        <Route path="dashboard" element={<Index />} />
      </Route>
    </Routes>
  </BrowserRouter>
);
