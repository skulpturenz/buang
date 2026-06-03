import { yupResolver } from "@hookform/resolvers/yup";
import { useForm } from "react-hook-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useSearchParams } from "react-router-dom";
import * as yup from "yup";
import { patch } from "../lib/api.js";

const schema = yup.object({
	buangApiBaseUrl: yup
		.string()
		.url("Must be a valid URL")
		.required("Required"),
	buangApiKey: yup.string().required("Required"),
});

type FormValues = yup.InferType<typeof schema>;

export const CompleteProfile = () => {
	const navigate = useNavigate();
	const [params] = useSearchParams();
	const queryClient = useQueryClient();
	const {
		register,
		handleSubmit,
		formState: { errors },
	} = useForm<FormValues>({ resolver: yupResolver(schema) });

	const mutation = useMutation({
		mutationFn: (values: FormValues) => patch("/me", values),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["me"] });
			const repo = params.get("repo");
			const installationId = params.get("installation_id");
			if (repo && installationId) {
				navigate(
					`/setup?repo=${encodeURIComponent(repo)}&installation_id=${installationId}`,
				);
			} else {
				navigate("/setup");
			}
		},
	});

	const onSubmit = async (values: FormValues) => {
		await mutation.mutateAsync(values);
	};

	return (
		<div className="mx-auto mt-20 max-w-sm space-y-6 px-4">
			<div>
				<h1 className="text-2xl font-semibold">Almost there</h1>
				<p className="mt-1 text-sm text-muted-foreground">
					Connect your Buang server to enable preview deployments.
				</p>
			</div>

			<form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
				{mutation.error && (
					<p className="text-sm text-destructive">
						{mutation.error.message || "Failed to save. Please try again."}
					</p>
				)}

				<div className="space-y-1">
					<label className="text-sm font-medium">Buang API URL</label>
					<input
						type="url"
						placeholder="https://buang.example.com"
						{...register("buangApiBaseUrl")}
						className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
					/>
					{errors.buangApiBaseUrl && (
						<p className="text-xs text-destructive">
							{errors.buangApiBaseUrl.message}
						</p>
					)}
				</div>

				<div className="space-y-1">
					<label className="text-sm font-medium">Buang API Key</label>
					<input
						type="password"
						{...register("buangApiKey")}
						className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
					/>
					{errors.buangApiKey && (
						<p className="text-xs text-destructive">
							{errors.buangApiKey.message}
						</p>
					)}
				</div>

				<button
					type="submit"
					disabled={mutation.isPending}
					className="w-full rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:opacity-90 disabled:opacity-50">
					{mutation.isPending ? "Saving…" : "Save and continue"}
				</button>
			</form>
		</div>
	);
};