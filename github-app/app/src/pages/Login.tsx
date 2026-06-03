import { yupResolver } from "@hookform/resolvers/yup";
import { useForm } from "react-hook-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate } from "react-router-dom";
import * as yup from "yup";
import { post } from "../lib/api.js";

const schema = yup.object({
	email: yup.string().email("Invalid email").required("Required"),
	password: yup.string().required("Required"),
});

type FormValues = yup.InferType<typeof schema>;

export const Login = () => {
	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const {
		register,
		handleSubmit,
		formState: { errors },
	} = useForm<FormValues>({ resolver: yupResolver(schema) });

	const mutation = useMutation({
		mutationFn: (values: FormValues) => post("/login", values),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["me"] });
			navigate("/setup");
		},
	});

	const onSubmit = async (values: FormValues) => {
		await mutation.mutateAsync(values);
	};

	return (
		<div className="mx-auto mt-20 max-w-sm space-y-6 px-4">
			<h1 className="text-2xl font-semibold">Sign in to Buang</h1>

			<form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
				{mutation.error && (
					<p className="text-sm text-destructive">
						{mutation.error.message || "Invalid email or password"}
					</p>
				)}

				<div className="space-y-1">
					<label className="text-sm font-medium">Email</label>
					<input
						type="email"
						{...register("email")}
						className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
					/>
					{errors.email && (
						<p className="text-xs text-destructive">{errors.email.message}</p>
					)}
				</div>

				<div className="space-y-1">
					<label className="text-sm font-medium">Password</label>
					<input
						type="password"
						{...register("password")}
						className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
					/>
					{errors.password && (
						<p className="text-xs text-destructive">{errors.password.message}</p>
					)}
				</div>

				<button
					type="submit"
					disabled={mutation.isPending}
					className="w-full rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:opacity-90 disabled:opacity-50">
					{mutation.isPending ? "Signing in…" : "Sign in"}
				</button>
			</form>

			<div className="relative">
				<div className="absolute inset-0 flex items-center">
					<span className="w-full border-t border-border" />
				</div>
				<div className="relative flex justify-center text-xs uppercase">
					<span className="bg-background px-2 text-muted-foreground">or</span>
				</div>
			</div>

			<a
				href="/api/auth/oidc"
				className="block w-full rounded-md border border-input bg-background px-4 py-2 text-center text-sm font-medium hover:bg-accent">
				Sign in with SSO
			</a>

			<p className="text-center text-sm text-muted-foreground">
				No account?{" "}
				<Link to="/register" className="text-primary underline">
					Register
				</Link>
			</p>
		</div>
	);
};