/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable jsx-a11y/html-has-lang */

"use client";

import { ClerkProvider } from "@clerk/nextjs";

import "../styles/globals.css";

export default function RootLayout({ children }) {
  return (
    <ClerkProvider>
      <html className="bg-primary-black">
        <body>
          {/* Main Content */}
          <main className="flex-1  overflow-hidden">{children}</main>
        </body>
      </html>
    </ClerkProvider>
  );
}
