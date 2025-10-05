<script lang="ts">
	import { TrendingUp, TrendingDown, Plus, Minus, Calendar, Play, Loader2 } from 'lucide-svelte';
	import type { ChangeDetectionResult } from '$lib/types';
	
    interface Props {
        onAnalyze: (basePeriod: number, comparePeriod: number) => Promise<void>;
        result: ChangeDetectionResult | null;
        loading: boolean;
        yearsMode?: boolean;
        startYear?: number;
        endYear?: number;
        onYearsChange?: (startYear: number, endYear: number) => void;
    }
	
    let { onAnalyze, result, loading, yearsMode = false, startYear, endYear, onYearsChange }: Props = $props();
	
    let basePeriod = $state(7); // days ago (used when yearsMode === false)
    let comparePeriod = $state(1); // days ago (used when yearsMode === false)
	let isOpen = $state(false);
	
    async function handleAnalyze() {
        if (yearsMode) {
            // In years mode, the parent handles years; pass placeholders
            await onAnalyze(0, 0);
            return;
        }
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
                <!-- Years Mode or Days Mode -->
                {#if yearsMode}
                    <div class="space-y-4 mb-4">
                        <div class="grid grid-cols-2 gap-3">
                            <div>
                                <label for="start-year" class="block text-sm font-medium text-slate-700 mb-2">Before (year)</label>
                                <input id="start-year" type="number" value={startYear ?? ''} min="1980" max="2100"
                                    oninput={(e: any) => onYearsChange?.(Number(e.target.value), endYear ?? 0)}
                                    class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm" />
                            </div>
                            <div>
                                <label for="end-year" class="block text-sm font-medium text-slate-700 mb-2">After (year)</label>
                                <input id="end-year" type="number" value={endYear ?? ''} min="1980" max="2100"
                                    oninput={(e: any) => onYearsChange?.(startYear ?? 0, Number(e.target.value))}
                                    class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm" />
                            </div>
                        </div>
                        <button
                            onclick={handleAnalyze}
                            disabled={loading || (startYear ?? 0) >= (endYear ?? 0)}
                            class="w-full flex items-center justify-center space-x-2 px-4 py-2 bg-orange-600 text-white rounded-md text-sm font-medium hover:bg-orange-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {#if loading}
                                <Loader2 class="w-4 h-4 animate-spin" />
                                <span>Analyzing...</span>
                            {:else}
                                <Play class="w-4 h-4" />
                                <span>Analyze Viewport Changes</span>
                            {/if}
                        </button>
                    </div>
                {:else}
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
                {/if}
				
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
