import { yupResolver } from "@hookform/resolvers/yup";
import { useForm } from "react-hook-form";
import { useNavigate, useSearchParams } from "react-router-dom";
import * as yup from "yup";
import { patch } from "../lib/api.js";
import { queryClient } from "../lib/queryClient.js";

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
	const {
		register,
		handleSubmit,
		setError,
		formState: { errors, isSubmitting },
	} = useForm<FormValues>({ resolver: yupResolver(schema) });

	const onSubmit = async (values: FormValues) => {
		try {
			await patch("/me", values);
			await queryClient.invalidateQueries({ queryKey: ["me"] });
			const repo = params.get("repo");
			const installationId = params.get("installation_id");
			if (repo && installationId) {
				navigate(
					`/setup?repo=${encodeURIComponent(repo)}&installation_id=${installationId}`,
				);
			} else {
				navigate("/setup");
			}
		} catch (_err) {
			setError("root", { message: "Failed to save. Please try again." });
		}
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
				{errors.root && (
					<p className="text-sm text-destructive">{errors.root.message}</p>
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
					disabled={isSubmitting}
					className="w-full rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:opacity-90 disabled:opacity-50">
					{isSubmitting ? "Saving…" : "Save and continue"}
				</button>
			</form>
		</div>
	);
};
