import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "react-toastify";

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
				toast("error try again");
			}
		};

		checkAuth();
	}, []);

	useEffect(() => {
		if (isAuthenticated === false) {
			navigate("/auth");
		}
	}, [isAuthenticated, navigate]);

	if (isAuthenticated === null) {
		return <div>Loading...</div>;
	}

	return null;
}

export default CheckAuth;
