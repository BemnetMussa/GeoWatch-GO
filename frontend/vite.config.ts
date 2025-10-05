import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
// import cesium from "vite-plugin-cesium" // <-- REMOVE

export default defineConfig({
    plugins: [tailwindcss(), sveltekit()/*, cesium()*/] // <-- REMOVE
});