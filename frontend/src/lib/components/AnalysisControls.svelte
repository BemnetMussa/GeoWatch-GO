<!-- src/lib/components/AnalysisControls.svelte -->

<script lang="ts">
	import type { Writable } from 'svelte/store';

	export let onAnalyze: () => void;
	export let loading: boolean = false;
	export let isZoomedInEnough: boolean = false;
	export let currentAltitude: number = 0;
	export let viewportArea: number = 0;
</script>

<div class="absolute top-4 left-1/2 -translate-x-1/2 z-10 p-4 bg-slate-800 bg-opacity-90 rounded-lg shadow-lg text-center">
	<!-- Zoom Status Indicator -->
	<div class="mb-3 text-sm">
		<div class="flex items-center justify-center space-x-2 mb-2">
			<div class="w-3 h-3 rounded-full {isZoomedInEnough ? 'bg-green-500' : 'bg-red-500'}"></div>
			<span class="text-slate-300">
				{isZoomedInEnough ? 'Analysis Ready' : 'Zoom In Required'}
			</span>
		</div>
		<div class="text-xs text-slate-400">
			Altitude: {(currentAltitude / 1000).toFixed(1)}km | 
			Area: {viewportArea.toFixed(1)} km²
		</div>
	</div>

	<!-- Analysis Button -->
	<button 
		on:click={onAnalyze}
		disabled={loading || !isZoomedInEnough}
		class="px-4 py-2 text-white font-bold rounded-md transition-colors w-full
		       {loading ? 'bg-gray-500 cursor-not-allowed' : 
		         isZoomedInEnough ? 'bg-blue-600 hover:bg-blue-700' : 
		         'bg-gray-500 cursor-not-allowed'}"
		title={!isZoomedInEnough ? 'Please zoom in closer to enable analysis' : 'Analyze changes in the current viewport'}
	>
		{#if loading}
			<div class="animate-spin rounded-full h-5 w-5 border-b-2 border-white mx-auto"></div>
		{:else if !isZoomedInEnough}
			🔍 Zoom In to Analyze
		{:else}
			📊 Analyze Viewport Changes
		{/if}
	</button>

	<!-- Helpful Tip -->
	{#if !isZoomedInEnough}
		<div class="mt-2 text-xs text-slate-400 max-w-xs">
			💡 Zoom in until you can see individual buildings or features clearly for best results
		</div>
	{/if}
</div>