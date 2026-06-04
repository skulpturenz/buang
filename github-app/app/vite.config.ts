import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";
import { defineConfig } from "vite";

export default defineConfig({
	root: "app",
	plugins: [react(), tailwindcss()],
	resolve: {
		alias: { "@": resolve(__dirname, "src") },
	},
	build: {
		outDir: "dist",
		emptyOutDir: true,
	},
	server: {
		proxy: {
			"/api": "http://localhost:3000",
		},
	},
});
