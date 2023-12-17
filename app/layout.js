/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable jsx-a11y/html-has-lang */

"use client";

import { ClerkProvider } from "@clerk/nextjs";

import Sidebar from "../components/ui/sidebar";

import "../styles/globals.css";

export default function RootLayout({ children }) {
  return (
    <ClerkProvider>
      <html>
        <body>
          <div className="flex h-screen bg-[#1B1B27]">
            {/* Sidebar (Navbar) */}

            <Sidebar />

            {/* Main Content */}
            <main className="flex-1 p-4 overflow-hidden">{children}</main>
          </div>
        </body>
      </html>
    </ClerkProvider>
  );
}
