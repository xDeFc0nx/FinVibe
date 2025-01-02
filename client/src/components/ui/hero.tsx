import { Component } from "solid-js";
import { Button } from "./button";
import heroimg from "@/assets/heroimg.svg";
import { WavyBackground } from "./wavy-background";
import { A } from "@solidjs/router";

const hero: Component = () => {
  return (
    <WavyBackground class="max-w-6xl mx-auto pb-40">
      <div class="container py-24 lg:py-32">
        <div class="grid md:grid-cols-2 gap-4 md:gap-8 xl:gap-20 md:items-center">
          <div>
            <h1 class="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
              <p class=" lg:text-5xl md:text-lg sm:text-xs text-white font-bold">
                FinVibe <span class="text-primary-pink"> elevate </span>
                your Finances
              </p>
            </h1>
            <p class="mt-3 text-xl text-muted-foreground">
              Empower your finances with FinVibe's smart tools and expert
              guidance
            </p>

            <div class="mt-7 grid gap-3 w-full sm:inline-flex ">
              <Button as={A} href="/register" variant="secondary" size={"lg"}>
                Get started
              </Button>
            </div>

            <div class="mt-6 lg:mt-10 grid grid-cols-2 gap-x-5" />
          </div>

          <div class="relative ms-4">
            <img
              class="w-full rounded-md"
              src={heroimg}
              alt="Image Description"
            />
          </div>
        </div>
      </div>
    </WavyBackground>
  );
};

export default hero;
