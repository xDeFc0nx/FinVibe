import { createSignal, onMount, createEffect } from "solid-js";
import { useNavigate } from "@solidjs/router";

export function CheckAuth(props) {
  const [isAuthenticated, setIsAuthenticated] = createSignal(null);
  const navigate = useNavigate();

  const checkAuth = async () => {
    try {
      const response = await fetch("http://localhost:3001/checkAuth", {
        method: "GET",
        credentials: "include",
      });

      if (response.ok) {
        setIsAuthenticated(true);
      } else {
        setIsAuthenticated(false);
      }
    } catch (error) {
      setIsAuthenticated(false);
    }
  };

  checkAuth();

  createEffect(() => {
    if (isAuthenticated() === false) {
      navigate("/login");
    }
  });

  createEffect(() => {
    if (isAuthenticated() === null) {
      return <div>Loading...</div>;
    }
  });

  return <></>;
}

export default CheckAuth;
