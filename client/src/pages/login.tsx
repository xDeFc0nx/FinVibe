import { createSignal } from "solid-js";
import {
  TextField,
  TextFieldInput,
  TextFieldLabel,
} from "@/components/ui/text-field";
import { A, useNavigate } from "@solidjs/router";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
export default function Login() {
  const [data, setData] = createSignal({});

  const [email, setEmail] = createSignal({});
  const [password, setPassword] = createSignal({});

  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    console.log("Email:", email());
    console.log("Password:", password());
    const response = await fetch("http://localhost:3001/Login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: email(),
        password: password(),
      }),
      credentials: "include",
    });
    const responseData = await response.json();
    console.log(responseData);
    setData(responseData);
    if (response.ok) {
      document.cookie = `jwt=${responseData.token}; path=/; secure; httpOnly`;

      navigate("/app/dashboard");
    } else {
      console.error("Login failed:", responseData.error);
    }
  };

  return (
    <>
      <div class="container relative hidden h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-1 lg:px-0">
        <Button
          as={A}
          href="/register"
          class={cn(
            buttonVariants({ variant: "ghost" }),
            "absolute right-4 top-4 md:right-8 md:top-8"
          )}
        >
          Register
        </Button>

        <div class="lg:p-8">
          <div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <div class="flex flex-col space-y-2 text-center">
              <h1 class="text-2xl font-semibold tracking-tight">
                {" "}
                Welcome Back!
              </h1>
              <p class="text-sm text-muted-foreground">
                Your data has been carefully wrapped in a secure layer of
                encryption and stored in a digital vault,
              </p>
            </div>

            <form action="submit" onSubmit={handleSubmit}>
              <TextField class="grid w-full max-w-sm items-center gap-1.5">
                <TextFieldLabel for="email">Email</TextFieldLabel>
                <TextFieldInput
                  required
                  type="email"
                  name="Email"
                  onChange={(e) => setEmail(e.target.value)}
                />
              </TextField>

              <TextField class="grid w-full max-w-sm items-center gap-1.5 pt-5">
                <TextFieldLabel for="Password">Password</TextFieldLabel>
                <TextFieldInput
                  required
                  type="password"
                  name="Password"
                  onChange={(e) => setPassword(e.target.value)}
                />
                <Button type="submit" variant="secondary">
                  Login
                </Button>
              </TextField>
            </form>
          </div>

          <p class="px-8 text-center text-sm text-muted-foreground pt-5">
            By clicking continue, you agree to our{" "}
            <a
              href="/terms"
              class="underline underline-offset-4 hover:text-primary"
            >
              Terms of Service
            </a>{" "}
            and{" "}
            <a
              href="/privacy"
              class="underline underline-offset-4 hover:text-primary"
            >
              Privacy Policy
            </a>
            .
          </p>
        </div>
      </div>
    </>
  );
}
