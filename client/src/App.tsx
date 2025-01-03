import { Button } from "@/components/ui/button";

import { ToastContainer } from "react-toastify";

import { Moon, Sun } from "lucide-react";
import React from "react";

import Hero from "@/components/ui/hero";
import { ThemeChanger } from "./components/ui/theme";

export default function App() {
  const [darkMode, setDarkMode] = React.useState(() =>
    document.documentElement.classList.contains("dark")
  );

  const toggleTheme = () => {
    setDarkMode(!darkMode);
    document.documentElement.classList.toggle("dark");
    localStorage.setItem("theme", !darkMode ? "dark" : "light");
  };

  React.useEffect(() => {
    const isDark =
      localStorage.getItem("theme") === "dark" ||
      (!localStorage.getItem("theme") &&
        window.matchMedia("(prefers-color-scheme: dark)").matches);

    setDarkMode(isDark);
    document.documentElement.classList.toggle("dark", isDark);
  }, []);

  return (
    <>
      <Hero />
      <ThemeChanger />
    </>
  );
}
