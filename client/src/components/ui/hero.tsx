import heroimg from "@/assets/heroimg.svg";
import type { Component } from "react";
import { Link } from "react-router";
import { Button } from "./button";
import { WavyBackground } from "./wavy-background";
const hero = () => {
  return (
    <>
      <WavyBackground classNameName="max-w-6xl mx-auto pt-64">
        <div className="container py-24 lg:py-32">
          <div className="grid md:grid-cols-2 gap-4  md:gap-8 xl:gap-20 md:items-center">
            <div>
              <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
                <p className=" lg:text-5xl md:text-lg sm:text-xs text-white font-bold">
                  FinVibe <span className="text-primary-pink"> elevate </span>
                  your Finances
                </p>
              </h1>
              <p className="mt-3 text-xl text-muted-foreground">
                Empower your finances with FinVibe's smart tools and expert
                guidance
              </p>

              <div className="mt-7 grid gap-3 w-full sm:inline-flex ">
                <Link to={"/login"}>
                  <Button>Get started</Button>
                </Link>
              </div>

              <div className="mt-6 lg:mt-10 grid grid-cols-2 gap-x-5" />
            </div>

            <div className="relative ms-4">
              <img
                className="w-full rounded-md"
                src={heroimg}
                alt="financial svg"
              />
            </div>
          </div>
        </div>
      </WavyBackground>
    </>
  );
};

export default hero;
