import { yupResolver } from "@hookform/resolvers/yup";
import { useForm } from "react-hook-form";
import { useSearchParams } from "react-router-dom";
import * as yup from "yup";
import { EnvVarsTable } from "../components/EnvVarsTable.js";
import { useSetupRepo } from "../hooks/useSetupRepo.js";

const schema = yup.object({
	repoFullName: yup.string().required("Required"),
	installationId: yup.number().required(),
	composePath: yup.string().required("Required"),
	serviceEntrypoint: yup.string().required("Required"),
	requiresAuthn: yup.boolean().required().default(false),
	username: yup.string().when("requiresAuthn", {
		is: true,
		then: s => s.required("Required when auth is enabled"),
		otherwise: s => s.optional(),
	}),
	password: yup.string().when("requiresAuthn", {
		is: true,
		then: s => s.required("Required — use a Personal Access Token (PAT), not your real password"),
		otherwise: s => s.optional(),
	}),
	envVars: yup
		.array(
			yup.object({ key: yup.string().required(), value: yup.string().required() }),
		)
		.default([]),
	waitForWorkflow: yup.boolean().default(false),
	workflowName: yup.string().nullable(),
});

type FormValues = yup.InferType<typeof schema>;

export const Setup = () => {
	const [params] = useSearchParams();
	const repoParam = params.get("repo") ?? "";
	const installationParam = params.get("installation_id") ?? "";

	const { mutateAsync, isPending, isSuccess, error } = useSetupRepo();

	const {
		register,
		control,
		watch,
		handleSubmit,
		formState: { errors },
	} = useForm<FormValues>({
		resolver: yupResolver(schema),
		defaultValues: {
			repoFullName: repoParam,
			installationId: installationParam ? Number(installationParam) : undefined,
			requiresAuthn: false,
			waitForWorkflow: false,
			envVars: [],
		},
	});

	const requiresAuthn = watch("requiresAuthn");
	const waitForWorkflow = watch("waitForWorkflow");

	const onSubmit = async (values: FormValues) => {
		await mutateAsync({
			...values,
			envVars: values.envVars ?? [],
			workflowName: values.workflowName ?? undefined,
		});
	};

	const apiErr =
		(error as any)?.body?.errors?.[0] ?? (error ? "Setup failed" : null);

	if (isSuccess) {
		return (
			<div className="mx-auto mt-20 max-w-sm px-4 text-center space-y-2">
				<p className="text-lg font-medium">Repository configured!</p>
				<p className="text-sm text-muted-foreground">
					Buang will now automatically deploy previews when pull requests are
					opened.
				</p>
			</div>
		);
	}

	return (
		<div className="mx-auto mt-10 max-w-lg space-y-6 px-4 pb-16">
			<div>
				<h1 className="text-2xl font-semibold">Configure repository</h1>
				<p className="mt-1 text-sm text-muted-foreground">
					Set up Buang preview deployments for{" "}
					<span className="font-mono font-medium">{repoParam}</span>.
				</p>
			</div>

			<form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
				{apiErr && <p className="text-sm text-destructive">{apiErr}</p>}

				<input type="hidden" {...register("repoFullName")} />
				<input type="hidden" {...register("installationId")} />

				<div className="space-y-1">
					<label className="text-sm font-medium">Compose file path</label>
					<input
						placeholder="docker-compose.yml"
						{...register("composePath")}
						className="w-full rounded-md border border-input bg-background px-3 py-2 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-ring"
					/>
					{errors.composePath && (
						<p className="text-xs text-destructive">{errors.composePath.message}</p>
					)}
				</div>

				<div className="space-y-1">
					<label className="text-sm font-medium">Service entrypoint</label>
					<input
						placeholder="traefik:80"
						{...register("serviceEntrypoint")}
						className="w-full rounded-md border border-input bg-background px-3 py-2 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-ring"
					/>
					<p className="text-xs text-muted-foreground">
						The service that accepts HTTP traffic (e.g. traefik:80, nginx:80).
					</p>
					{errors.serviceEntrypoint && (
						<p className="text-xs text-destructive">
							{errors.serviceEntrypoint.message}
						</p>
					)}
				</div>

				<div className="flex items-center gap-3">
					<input
						type="checkbox"
						id="requiresAuthn"
						{...register("requiresAuthn")}
						className="h-4 w-4 rounded border-input"
					/>
					<label htmlFor="requiresAuthn" className="text-sm font-medium">
						Repository requires authentication
					</label>
				</div>

				{requiresAuthn && (
					<div className="space-y-4 rounded-md border border-input p-4">
						<div className="space-y-1">
							<label className="text-sm font-medium">Username</label>
							<input
								{...register("username")}
								className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
							/>
							{errors.username && (
								<p className="text-xs text-destructive">{errors.username.message}</p>
							)}
						</div>
						<div className="space-y-1">
							<label className="text-sm font-medium">Password / PAT</label>
							<input
								type="password"
								{...register("password")}
								className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
							/>
							<p className="text-xs text-muted-foreground">
								Use a Personal Access Token (PAT) instead of your real password.
							</p>
							{errors.password && (
								<p className="text-xs text-destructive">{errors.password.message}</p>
							)}
						</div>
					</div>
				)}

				<div className="space-y-2">
					<label className="text-sm font-medium">
						Environment variables{" "}
						<span className="text-muted-foreground">(optional)</span>
					</label>
					<EnvVarsTable control={control} name="envVars" />
				</div>

				<div className="flex items-center gap-3">
					<input
						type="checkbox"
						id="waitForWorkflow"
						{...register("waitForWorkflow")}
						className="h-4 w-4 rounded border-input"
					/>
					<label htmlFor="waitForWorkflow" className="text-sm font-medium">
						Wait for workflow run before deploying
					</label>
				</div>

				{waitForWorkflow && (
					<div className="space-y-1">
						<label className="text-sm font-medium">Workflow name</label>
						<input
							placeholder="Build and push"
							{...register("workflowName")}
							className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
						/>
						<p className="text-xs text-muted-foreground">
							Leave blank to trigger on any successful workflow run.
						</p>
					</div>
				)}

				<button
					type="submit"
					disabled={isPending}
					className="w-full rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:opacity-90 disabled:opacity-50">
					{isPending ? "Configuring…" : "Save configuration"}
				</button>
			</form>
		</div>
	);
};
