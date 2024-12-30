import { createEffect, createSignal, onCleanup, Suspense } from "solid-js";
import AboutData from "./about.data";
import { WebSocketClient } from "@/libs/socket";
import { useLocation } from "@solidjs/router";

export default function About() {
  const name = AboutData();

  createEffect(() => {
    console.log(name());
  });

  return (
    <section class="bg-pink-100 text-gray-700 p-8">
      <h1 class="text-2xl font-bold">About</h1>
    </section>
  );
}
