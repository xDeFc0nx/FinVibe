import { createEffect, createSignal, Suspense } from "solid-js";

import "@/index.css";
export default function Index() {
  const [data, setData] = createSignal({});

  const [email, setEmail] = createSignal({});
  const [password, setPassword] = createSignal({});

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
      credentials: "include", // Include cookies in the request
    });
    const responseData = await response.json();
    console.log(responseData);
    setData(responseData);
    if (response.ok) {
      console.log("Token:", responseData.token);

      document.cookie = `jwt=${responseData.token}; path=/; secure; httpOnly`;
    } else {
      console.error("Login failed:", responseData.error);
    }
  };

  return (
    <div class="w-full h-full bg-black">
      <form action="submit" onsubmit={handleSubmit}>
        <input
          type="Email"
          name="Email"
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="Password"
          name="password"
          onChange={(e) => setPassword(e.target.value)}
        />
        <button type="submit">Login</button>
      </form>
    </div>
  );
}
