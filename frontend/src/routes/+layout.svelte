<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
</script>

<div class="min-h-screen bg-slate-50">
	<Header />
	<main class="flex-1">
		{@render children?.()}
	</main>
</div>
