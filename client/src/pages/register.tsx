import { createSignal } from "solid-js";
import {
  TextField,
  TextFieldInput,
  TextFieldLabel,
} from "@/components/ui/text-field";
import { A, useNavigate } from "@solidjs/router";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import toast from "solid-toast";
export default function Register() {
  const [data, setData] = createSignal({});

  const [firstName, setFirstName] = createSignal("");
  const [lastName, setLastName] = createSignal("");
  const [country, setCountry] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const [emailError, setEmailError] = createSignal("");

  const navigate = useNavigate();

  const handlePasswordChange = (e) => {
    const newPassword = e.target.value;
    setPassword(newPassword);

    if (newPassword.length < 8) {
      setError("Password must be at least 8 characters long");
    } else {
      setError("");
    }
  };

  const handleEmailChange = (e) => {
    const newEmail = e.target.value;
    setEmail(newEmail);

    // Validate email format using a regex pattern
    const emailPattern = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailPattern.test(newEmail)) {
      setEmailError("Please enter a valid email address");
    } else {
      setEmailError("");
    }
  };
  const handleSubmit = async (event) => {
    event.preventDefault();

    const response = await fetch("http://localhost:3001/Register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        firstName: firstName(),
        lastName: lastName(),
        country: country(),
        email: email(),
        password: password(),
      }),
      credentials: "include",
    });
    const responseData = await response.json();
    console.log(responseData);
    setData(responseData);
    if (response.ok) {
      toast.success("Successfully Registered, Please Login");

      navigate("/login");
    } else {
      toast.error("Registration failed, please try again");
    }
  };

  return (
    <>
      <div class="container relative hidden h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0">
        <Button
          as={A}
          href="/login"
          class={cn(
            buttonVariants({ variant: "ghost" }),
            "absolute right-4 top-4 md:right-8 md:top-8"
          )}
        >
          Login
        </Button>
        <div class="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex">
          <div class="absolute inset-0 bg-zinc-900" />
          <div class="relative z-20 flex items-center text-lg font-medium">
            FinVibe
          </div>
          <div class="relative z-20 mt-auto">
            <blockquote class="space-y-2">
              <p class="text-lg">&ldquo;HI&rdquo;</p>
              <footer class="text-sm"> Nehar Tale</footer>
            </blockquote>
          </div>
        </div>
        <div class="lg:p-8">
          <div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <div class="flex flex-col space-y-2 text-center">
              <h1 class="text-2xl font-semibold tracking-tight"> Welcome!</h1>
              <p class="text-sm text-muted-foreground">
                Get ready for a smooth ride to financial freedom—no boring
                spreadsheets here. Let’s make money management a little more
                fun, shall we?
              </p>
            </div>

            <form action="submit" onSubmit={handleSubmit}>
              <TextField class="grid  max-w-sm items-center gap-1.5 pt-5">
                <TextFieldLabel for="email">First Name</TextFieldLabel>
                <TextFieldInput
                  required
                  type="text"
                  name="FirstName"
                  onChange={(e) => setFirstName(e.target.value)}
                />
              </TextField>
              <TextField class="grid w-full max-w-sm items-center gap-1.5 pt-5">
                <TextFieldLabel for="email">Last Name</TextFieldLabel>
                <TextFieldInput
                  required
                  type="text"
                  name="LastName"
                  onChange={(e) => setLastName(e.target.value)}
                />
              </TextField>

              <TextField class="grid w-full max-w-sm items-center gap-1.5 pt-5">
                <TextFieldLabel for="email">Country</TextFieldLabel>
                <TextFieldInput
                  required
                  type="text"
                  name="Country"
                  onChange={(e) => setCountry(e.target.value)}
                />
              </TextField>
              <TextField class="grid w-full max-w-sm items-center gap-1.5 pt-5">
                <TextFieldLabel for="email">Email</TextFieldLabel>
                <TextFieldInput
                  required
                  type="email"
                  name="Email"
                  value={email()}
                  onInput={handleEmailChange}
                />
                {emailError() && (
                  <div class="text-red-500 text-sm">{emailError()}</div>
                )}
              </TextField>
              <TextField class="grid w-full max-w-sm items-center gap-1.5 pt-5">
                <TextFieldLabel for="password">Password</TextFieldLabel>
                <TextFieldInput
                  required
                  type="password"
                  name="Password"
                  value={password()}
                  onInput={handlePasswordChange}
                />
                {error() && <div class="text-red-500 text-sm">{error()}</div>}
                <Button type="submit" variant="secondary">
                  Register
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
