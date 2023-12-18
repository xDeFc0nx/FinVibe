"use client";

import React from "react";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";

import Bg from "../components/ui/LandingPage/background";
import Hero from "../components/ui/LandingPage/hero";
import Cards from "../components/ui/LandingPage/Cards";

const page = () => (
  <>
    <ToastContainer />
    <Bg />
    <Hero />
    <Cards />
  </>
);

export default page;
