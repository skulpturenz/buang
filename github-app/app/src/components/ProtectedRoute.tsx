import { type ReactNode } from "react";
import { Navigate } from "react-router-dom";
import { useMe } from "../hooks/useMe.js";

export const ProtectedRoute = ({ children }: { children: ReactNode }) => {
	const { data, isLoading, isError } = useMe();

	if (isLoading) return <div className="p-8 text-center">Loading…</div>;
	if (isError || !data) return <Navigate to="/login" replace />;

	if (!data.hasBuangKey) return <Navigate to="/complete-profile" replace />;

	return <>{children}</>;
};
