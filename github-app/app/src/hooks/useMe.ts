import { useQuery } from "@tanstack/react-query";
import { get } from "../lib/api.js";

export interface Me {
	id: string;
	email: string | null;
	authType: "email" | "oauth";
	buangApiBaseUrl: string | null;
	hasBuangKey: boolean;
}

export const useMe = () =>
	useQuery<Me>({
		queryKey: ["me"],
		queryFn: () => get<Me>("/me"),
		retry: false,
	});
