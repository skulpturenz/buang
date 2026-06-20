import "../app/src/index.css";
import type { Preview } from "@storybook/react";

const preview: Preview = {
	parameters: {
		backgrounds: {
			default: "light",
			values: [
				{ name: "light", value: "oklch(0.985 0.003 110)" },
				{ name: "dark", value: "oklch(0.135 0.005 110)" },
			],
		},
	},
	decorators: [
		(Story, context) => {
			const isDark =
				context.globals.backgrounds?.value === "oklch(0.135 0.005 110)";
			document.body.className = isDark ? "theme-dark" : "theme-light";
			return Story();
		},
	],
};

export default preview;
