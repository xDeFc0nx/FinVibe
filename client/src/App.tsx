import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Toaster } from "@/components/ui/toaster";
import { useToast } from "@/hooks/use-toast";
import {
  BookOpen,
  Boxes,
  Cpu,
  FileCode2,
  Flame,
  GitFork,
  Github,
  Moon,
  Package,
  Palette,
  Star,
  Sun,
  Terminal,
  Wrench,
} from "lucide-react";
import React from "react";

import Hero from "@/components/ui/hero";

export default function App() {
  const { toast } = useToast();
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
      <div className="absolute top-4 right-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleTheme}
          className="rounded-full"
        >
          {darkMode ? <Sun className="size-5" /> : <Moon className="size-5" />}
        </Button>
      </div>

      <Toaster />
    </>
  );
}
