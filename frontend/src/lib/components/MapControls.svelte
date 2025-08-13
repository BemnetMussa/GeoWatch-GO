<script lang="ts">
	import { Layers, Flame, RefreshCw, Satellite, Map, Navigation, TrendingUp, BarChart3 } from 'lucide-svelte';
	
	interface Props {
		activeLayer: string;
		showFires: boolean;
		loading: boolean;
		fireCount: number;
		showChanges?: boolean;
		changeCount?: number;
		onLayerChange: (layer: 'satellite' | 'terrain' | 'street') => void;
		onToggleFires: () => void;
		onToggleChanges?: () => void;
		onRefresh: () => void;
		onShowAnalytics?: () => void;
	}
	
	let { 
		activeLayer, 
		showFires, 
		loading, 
		fireCount, 
		showChanges = false,
		changeCount = 0,
		onLayerChange, 
		onToggleFires, 
		onToggleChanges,
		onRefresh,
		onShowAnalytics
	}: Props = $props();
	
	let showLayerMenu = $state(false);
</script>

<!-- Layer Control -->
<div class="absolute top-4 right-4 z-10">
	<div class="bg-white rounded-lg shadow-lg border border-slate-200 overflow-hidden">
		<button
			onclick={() => showLayerMenu = !showLayerMenu}
			class="flex items-center space-x-2 px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors w-full"
		>
			<Layers class="w-4 h-4" />
			<span class="capitalize">{activeLayer}</span>
		</button>
		
		{#if showLayerMenu}
			<div class="border-t border-slate-200">
				<button
					onclick={() => { onLayerChange('satellite'); showLayerMenu = false; }}
					class="flex items-center space-x-2 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 transition-colors w-full {activeLayer === 'satellite' ? 'bg-orange-50 text-orange-700' : ''}"
				>
					<Satellite class="w-4 h-4" />
					<span>Satellite</span>
				</button>
				<button
					onclick={() => { onLayerChange('terrain'); showLayerMenu = false; }}
					class="flex items-center space-x-2 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 transition-colors w-full {activeLayer === 'terrain' ? 'bg-orange-50 text-orange-700' : ''}"
				>
					<Navigation class="w-4 h-4" />
					<span>Terrain</span>
				</button>
				<button
					onclick={() => { onLayerChange('street'); showLayerMenu = false; }}
					class="flex items-center space-x-2 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 transition-colors w-full {activeLayer === 'street' ? 'bg-orange-50 text-orange-700' : ''}"
				>
					<Map class="w-4 h-4" />
					<span>Street</span>
				</button>
			</div>
		{/if}
	</div>
</div>

<!-- Fire Toggle & Controls -->
<div class="absolute top-4 left-4 z-10 flex flex-col space-y-2">
	<button
		onclick={onToggleFires}
		class="flex items-center space-x-2 px-4 py-3 bg-white rounded-lg shadow-lg border border-slate-200 text-sm font-medium transition-colors {showFires ? 'text-orange-700 bg-orange-50' : 'text-slate-700 hover:bg-slate-50'}"
	>
		<Flame class="w-4 h-4" />
		<div class="flex flex-col items-start">
			<span>Fire Data</span>
			{#if fireCount > 0}
				<span class="text-xs text-slate-500">{fireCount} active</span>
			{/if}
		</div>
	</button>
	
	{#if onToggleChanges}
		<button
			onclick={onToggleChanges}
			class="flex items-center space-x-2 px-4 py-3 bg-white rounded-lg shadow-lg border border-slate-200 text-sm font-medium transition-colors {showChanges ? 'text-blue-700 bg-blue-50' : 'text-slate-700 hover:bg-slate-50'}"
		>
			<TrendingUp class="w-4 h-4" />
			<div class="flex flex-col items-start">
				<span>Changes</span>
				{#if changeCount > 0}
					<span class="text-xs text-slate-500">{changeCount} new</span>
				{/if}
			</div>
		</button>
	{/if}
	
	{#if onShowAnalytics}
		<button
			onclick={onShowAnalytics}
			class="flex items-center space-x-2 px-4 py-3 bg-white rounded-lg shadow-lg border border-slate-200 text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors"
		>
			<BarChart3 class="w-4 h-4" />
			<span>Analytics</span>
		</button>
	{/if}
	
	<button
		onclick={onRefresh}
		disabled={loading}
		class="flex items-center space-x-2 px-4 py-3 bg-white rounded-lg shadow-lg border border-slate-200 text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors disabled:opacity-50"
	>
		<RefreshCw class="w-4 h-4 {loading ? 'animate-spin' : ''}" />
		<span>Refresh</span>
	</button>
</div>
