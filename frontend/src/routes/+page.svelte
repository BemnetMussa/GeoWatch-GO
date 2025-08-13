<script lang="ts">
	import MapContainer from '$lib/components/MapContainer.svelte';
	import { onMount } from 'svelte';
	
	let mounted = $state(false);
	
	onMount(() => {
		mounted = true;
	});
</script>

<svelte:head>
	<title>Wildfire Detection & Change Analysis</title>
	<meta name="description" content="Real-time wildfire detection and geographical change analysis using satellite data" />
</svelte:head>

<div class="h-screen flex flex-col">
	<div class="flex-1 relative">
		{#if mounted}
			<MapContainer />
		{:else}
			<div class="h-full flex items-center justify-center bg-slate-100">
				<div class="text-center">
					<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-600 mx-auto mb-4"></div>
					<p class="text-slate-600">Loading map...</p>
				</div>
			</div>
		{/if}
	</div>
</div>
