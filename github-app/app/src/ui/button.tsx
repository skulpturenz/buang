import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import * as React from "react";
import { cn } from "@/lib/utils";

const buttonVariants = cva(
	[
		"inline-flex items-center justify-center gap-2 w-full",
		"px-[14px] py-[10px]",
		"border border-[var(--line-strong)]",
		"bg-[var(--fg)] text-[var(--bg)]",
		"font-[family-name:var(--font-mono)] text-[12px] font-semibold uppercase tracking-[0.08em]",
		"rounded-none whitespace-nowrap cursor-pointer",
		"transition-[background,color] duration-[80ms] ease-linear",
		"hover:bg-[var(--bg)] hover:text-[var(--fg)]",
		"disabled:pointer-events-none disabled:opacity-50",
		"focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--line-strong)]",
	],
	{
		variants: {
			variant: {
				default: "",
				accent: [
					"bg-[var(--ds-accent)] border-[var(--ds-accent)] text-[var(--ds-accent-fg)]",
					"hover:bg-[var(--bg)] hover:text-[var(--fg)] hover:border-[var(--line-strong)]",
				],
				ghost: [
					"bg-transparent border-transparent text-[var(--fg)]",
					"hover:bg-[var(--fg)] hover:text-[var(--bg)]",
				],
				destructive: [
					"bg-[var(--ds-destructive)] border-[var(--ds-destructive)] text-[var(--ds-destructive-fg)]",
					"hover:bg-[var(--bg)] hover:text-[var(--fg)] hover:border-[var(--line-strong)]",
				],
			},
			size: {
				default: "",
				sm: "px-[10px] py-[6px] text-[10.5px]",
				icon: "w-auto px-[6px] py-[6px]",
			},
		},
		defaultVariants: { variant: "default", size: "default" },
	},
);

export interface ButtonProps
	extends React.ButtonHTMLAttributes<HTMLButtonElement>,
		VariantProps<typeof buttonVariants> {
	asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
	({ className, variant, size, asChild = false, ...props }, ref) => {
		const Comp = asChild ? Slot : "button";
		return (
			<Comp
				className={cn(buttonVariants({ variant, size, className }))}
				ref={ref}
				{...props}
			/>
		);
	},
);
Button.displayName = "Button";

export { Button, buttonVariants };
