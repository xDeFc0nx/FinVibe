import { createEffect, createSignal, onCleanup, Suspense } from "solid-js";
import { WebSocketClient } from "@/lib/socket";
import { useLocation } from "@solidjs/router";

export default function About() {
  return (
    <section class="bg-pink-100 text-gray-700 p-8">
      <h1 class="text-2xl font-bold">About</h1>
    </section>
  );
}
