<script lang="ts">
	import { TrendingUp, TrendingDown, Plus, Minus, Calendar, Play, Loader2 } from 'lucide-svelte';
	import type { ChangeDetectionResult } from '$lib/types';
	
	interface Props {
		onAnalyze: (basePeriod: number, comparePeriod: number) => Promise<void>;
		result: ChangeDetectionResult | null;
		loading: boolean;
	}
	
	let { onAnalyze, result, loading }: Props = $props();
	
	let basePeriod = $state(7); // 1 week ago
	let comparePeriod = $state(1); // 1 day ago
	let isOpen = $state(false);
	
	async function handleAnalyze() {
		await onAnalyze(basePeriod, comparePeriod);
	}
	
	function getChangeColor(type: 'new' | 'extinguished' | 'growing' | 'diminishing') {
		switch (type) {
			case 'new': return 'text-red-600 bg-red-50';
			case 'extinguished': return 'text-green-600 bg-green-50';
			case 'growing': return 'text-orange-600 bg-orange-50';
			case 'diminishing': return 'text-blue-600 bg-blue-50';
		}
	}
</script>

<div class="absolute top-20 right-4 z-10">
	<div class="bg-white rounded-lg shadow-lg border border-slate-200 overflow-hidden">
		<button
			onclick={() => isOpen = !isOpen}
			class="flex items-center space-x-2 px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors w-full"
		>
			<TrendingUp class="w-4 h-4" />
			<span>Change Detection</span>
		</button>
		
		{#if isOpen}
			<div class="border-t border-slate-200 p-4 w-80">
				<!-- Time Period Controls -->
				<div class="space-y-4 mb-4">
					<div>
						<label for="basePeriod" class="block text-sm font-medium text-slate-700 mb-2">
							Base Period (days ago)
						</label>
						<select id="basePeriod" bind:value={basePeriod} class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm">
							<option value={1}>1 day ago</option>
							<option value={3}>3 days ago</option>
							<option value={7}>1 week ago</option>
							<option value={14}>2 weeks ago</option>
							<option value={30}>1 month ago</option>
						</select>
					</div>
					
					<div>
						<label for="comparePeriod" class="block text-sm font-medium text-slate-700 mb-2">
							Compare Period (days ago)
						</label>
						<select id="comparePeriod" bind:value={comparePeriod} class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm">
							<option value={0}>Current</option>
							<option value={1}>1 day ago</option>
							<option value={3}>3 days ago</option>
							<option value={7}>1 week ago</option>
						</select>
					</div>
					
					<button
						onclick={handleAnalyze}
						disabled={loading || basePeriod <= comparePeriod}
						class="w-full flex items-center justify-center space-x-2 px-4 py-2 bg-orange-600 text-white rounded-md text-sm font-medium hover:bg-orange-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{#if loading}
							<Loader2 class="w-4 h-4 animate-spin" />
							<span>Analyzing...</span>
						{:else}
							<Play class="w-4 h-4" />
							<span>Analyze Changes</span>
						{/if}
					</button>
				</div>
				
				<!-- Results -->
				{#if result}
					<div class="space-y-3">
						<div class="text-sm font-medium text-slate-900 border-b border-slate-200 pb-2">
							Change Summary ({result.summary.changePercentage}% change)
						</div>
						
						<!-- New Fires -->
						<div class="flex items-center justify-between p-2 rounded-lg {getChangeColor('new')}">
							<div class="flex items-center space-x-2">
								<Plus class="w-4 h-4" />
								<span class="text-sm font-medium">New Fires</span>
							</div>
							<span class="text-sm font-bold">{result.summary.totalNew}</span>
						</div>
						
						<!-- Extinguished Fires -->
						<div class="flex items-center justify-between p-2 rounded-lg {getChangeColor('extinguished')}">
							<div class="flex items-center space-x-2">
								<Minus class="w-4 h-4" />
								<span class="text-sm font-medium">Extinguished</span>
							</div>
							<span class="text-sm font-bold">{result.summary.totalExtinguished}</span>
						</div>
						
						<!-- Growing Fires -->
						<div class="flex items-center justify-between p-2 rounded-lg {getChangeColor('growing')}">
							<div class="flex items-center space-x-2">
								<TrendingUp class="w-4 h-4" />
								<span class="text-sm font-medium">Growing</span>
							</div>
							<span class="text-sm font-bold">{result.summary.totalGrowing}</span>
						</div>
						
						<!-- Diminishing Fires -->
						<div class="flex items-center justify-between p-2 rounded-lg {getChangeColor('diminishing')}">
							<div class="flex items-center space-x-2">
								<TrendingDown class="w-4 h-4" />
								<span class="text-sm font-medium">Diminishing</span>
							</div>
							<span class="text-sm font-bold">{result.summary.totalDiminishing}</span>
						</div>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>
