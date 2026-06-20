import { Trash2 } from "lucide-react";
import {
	type Control,
	useFieldArray,
} from "react-hook-form";

interface Props {
	control: Control<any>;
	name: string;
}

export const EnvVarsTable = ({ control, name }: Props) => {
	const { fields, append, remove } = useFieldArray({ control, name });

	return (
		<div className="space-y-2">
			<table className="w-full text-sm">
				<thead>
					<tr className="text-left text-muted-foreground">
						<th className="pb-1 font-medium">Key</th>
						<th className="pb-1 font-medium">Value</th>
						<th />
					</tr>
				</thead>
				<tbody className="space-y-1">
					{fields.map((field, index) => (
						<tr key={field.id}>
							<td className="pr-2">
								<input
									{...control.register(`${name}.${index}.key`)}
									placeholder="KEY"
									className="w-full rounded border border-input bg-background px-2 py-1 font-mono text-xs focus:outline-none focus:ring-1 focus:ring-ring"
								/>
							</td>
							<td className="pr-2">
								<input
									{...control.register(`${name}.${index}.value`)}
									placeholder="value"
									className="w-full rounded border border-input bg-background px-2 py-1 font-mono text-xs focus:outline-none focus:ring-1 focus:ring-ring"
								/>
							</td>
							<td>
								<button
									type="button"
									onClick={() => remove(index)}
									className="text-destructive hover:opacity-70">
									<Trash2 size={14} />
								</button>
							</td>
						</tr>
					))}
				</tbody>
			</table>
			<button
				type="button"
				onClick={() => append({ key: "", value: "" })}
				className="text-xs text-primary underline underline-offset-2 hover:opacity-70">
				+ Add variable
			</button>
		</div>
	);
};
