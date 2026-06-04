import type { Meta, StoryObj } from "@storybook/react";
import { Search } from "lucide-react";
import { Button } from "@/ui/button";

const meta: Meta<typeof Button> = {
	title: "UI/Button",
	component: Button,
	tags: ["autodocs"],
	parameters: {
		layout: "centered",
	},
	argTypes: {
		variant: {
			control: "select",
			options: ["default", "accent", "ghost", "destructive"],
		},
		size: {
			control: "select",
			options: ["default", "sm", "icon"],
		},
		disabled: { control: "boolean" },
		asChild: { control: "boolean" },
	},
};

export default meta;
type Story = StoryObj<typeof Button>;

export const Default: Story = {
	args: { children: "Default", variant: "default" },
};

export const Accent: Story = {
	args: { children: "Accent", variant: "accent" },
};

export const Ghost: Story = {
	args: { children: "Ghost", variant: "ghost" },
};

export const Destructive: Story = {
	args: { children: "Destructive", variant: "destructive" },
};

export const Small: Story = {
	args: { children: "Small", size: "sm" },
};

export const Icon: Story = {
	args: {
		size: "icon",
		"aria-label": "Search",
		children: <Search size={14} />,
	},
};

export const Disabled: Story = {
	args: { children: "Disabled", disabled: true },
};

export const AllVariants: Story = {
	name: "All Variants",
	render: () => (
		<div className="flex flex-col gap-3" style={{ width: 280 }}>
			{(["default", "accent", "ghost", "destructive"] as const).map(
				(variant) => (
					<div key={variant} className="flex gap-2 items-center">
						<Button variant={variant} size="sm">
							{variant} sm
						</Button>
						<Button variant={variant}>{variant}</Button>
						<Button variant={variant} size="icon" aria-label={variant}>
							<Search size={14} />
						</Button>
					</div>
				),
			)}
		</div>
	),
};

export const AsChildLink: Story = {
	name: "asChild (renders as <a>)",
	render: () => (
		<Button asChild variant="accent">
			<a href="#" onClick={(e) => e.preventDefault()}>
				Link as button
			</a>
		</Button>
	),
};
