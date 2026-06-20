import tailwindcss from "@tailwindcss/vite";
import type { StorybookConfig } from "@storybook/react-vite";
import { resolve } from "path";

const config: StorybookConfig = {
	stories: ["../app/src/**/*.stories.@(ts|tsx)"],
	addons: ["@storybook/addon-essentials"],
	framework: { name: "@storybook/react-vite", options: {} },
	viteFinal(config) {
		config.plugins = [tailwindcss(), ...(config.plugins ?? [])];
		config.resolve = {
			...config.resolve,
			alias: {
				...((config.resolve?.alias as Record<string, string>) ?? {}),
				"@": resolve(__dirname, "../app/src"),
			},
		};
		return config;
	},
};

export default config;
