import { Toast } from "@/components/ui/toast";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

export function CheckAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState<null | boolean>(null);
  const navigate = useNavigate();

  useEffect(() => {
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
        Toast({
          title: "Error",
          variant: "destructive",
        });
      }
    };

    checkAuth();
  }, []);

  useEffect(() => {
    if (isAuthenticated === false) {
      navigate("/login");
    }
  }, [isAuthenticated, navigate]);

  if (isAuthenticated === null) {
    return <div>Loading...</div>;
  }

  return null;
}

export default CheckAuth;
