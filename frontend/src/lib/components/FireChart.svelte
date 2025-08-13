<script lang="ts">
	import { TrendingUp, Clock, BarChart3 } from 'lucide-svelte';
	import type { FireData } from '$lib/types';
	
	interface Props {
		fires: FireData[];
	}
	
	let { fires }: Props = $props();
	
	// Process data for charts
	let chartData = $derived.by(() => {
		if (!fires.length) return null;
		
		// Time distribution (by hour)
		const hourlyData = Array.from({ length: 24 }, (_, i) => ({
			hour: i.toString().padStart(2, '0'),
			count: 0,
			label: `${i.toString().padStart(2, '0')}:00`
		}));
		
		fires.forEach(fire => {
			const hour = parseInt(fire.acq_time.substring(0, 2));
			if (hour >= 0 && hour < 24) {
				hourlyData[hour].count++;
			}
		});
		
		// Brightness distribution
		const brightnessRanges = [
			{ min: 0, max: 300, label: '0-300K', count: 0 },
			{ min: 300, max: 350, label: '300-350K', count: 0 },
			{ min: 350, max: 400, label: '350-400K', count: 0 },
			{ min: 400, max: 450, label: '400-450K', count: 0 },
			{ min: 450, max: Infinity, label: '450K+', count: 0 }
		];
		
		fires.forEach(fire => {
			const range = brightnessRanges.find(r => fire.brightness >= r.min && fire.brightness < r.max);
			if (range) range.count++;
		});
		
		// FRP distribution
		const frpRanges = [
			{ min: 0, max: 10, label: '0-10 MW', count: 0 },
			{ min: 10, max: 50, label: '10-50 MW', count: 0 },
			{ min: 50, max: 100, label: '50-100 MW', count: 0 },
			{ min: 100, max: 500, label: '100-500 MW', count: 0 },
			{ min: 500, max: Infinity, label: '500+ MW', count: 0 }
		];
		
		fires.forEach(fire => {
			const frp = fire.frp || 0;
			const range = frpRanges.find(r => frp >= r.min && frp < r.max);
			if (range) range.count++;
		});
		
		const maxHourly = Math.max(...hourlyData.map(d => d.count));
		const maxBrightness = Math.max(...brightnessRanges.map(d => d.count));
		const maxFRP = Math.max(...frpRanges.map(d => d.count));
		
		return {
			hourlyData,
			brightnessRanges,
			frpRanges,
			maxHourly,
			maxBrightness,
			maxFRP
		};
	});
</script>

{#if chartData}
	<div class="space-y-8">
		<!-- Hourly Distribution -->
		<div class="bg-white rounded-lg border border-slate-200 p-6">
			<div class="flex items-center space-x-2 mb-6">
				<Clock class="w-5 h-5 text-orange-600" />
				<h3 class="text-lg font-semibold text-slate-900">Fire Detection by Hour</h3>
			</div>
			<div class="space-y-2">
				{#each chartData.hourlyData as data}
					<div class="flex items-center space-x-3">
						<span class="text-xs font-mono text-slate-500 w-12">{data.label}</span>
						<div class="flex-1 bg-slate-100 rounded-full h-4 relative">
							<div 
								class="bg-gradient-to-r from-orange-400 to-red-500 h-4 rounded-full transition-all duration-300"
								style="width: {chartData.maxHourly > 0 ? (data.count / chartData.maxHourly) * 100 : 0}%"
							></div>
							<span class="absolute right-2 top-0 text-xs font-medium text-slate-700 leading-4">
								{data.count}
							</span>
						</div>
					</div>
				{/each}
			</div>
		</div>
		
		<!-- Brightness Distribution -->
		<div class="bg-white rounded-lg border border-slate-200 p-6">
			<div class="flex items-center space-x-2 mb-6">
				<TrendingUp class="w-5 h-5 text-orange-600" />
				<h3 class="text-lg font-semibold text-slate-900">Brightness Distribution</h3>
			</div>
			<div class="space-y-3">
				{#each chartData.brightnessRanges as range}
					<div class="flex items-center space-x-3">
						<span class="text-sm text-slate-600 w-20">{range.label}</span>
						<div class="flex-1 bg-slate-100 rounded-full h-6 relative">
							<div 
								class="bg-gradient-to-r from-yellow-400 to-orange-500 h-6 rounded-full transition-all duration-300"
								style="width: {chartData.maxBrightness > 0 ? (range.count / chartData.maxBrightness) * 100 : 0}%"
							></div>
							<span class="absolute right-3 top-0 text-sm font-medium text-slate-700 leading-6">
								{range.count}
							</span>
						</div>
					</div>
				{/each}
			</div>
		</div>
		
		<!-- FRP Distribution -->
		<div class="bg-white rounded-lg border border-slate-200 p-6">
			<div class="flex items-center space-x-2 mb-6">
				<BarChart3 class="w-5 h-5 text-orange-600" />
				<h3 class="text-lg font-semibold text-slate-900">Fire Radiative Power Distribution</h3>
			</div>
			<div class="space-y-3">
				{#each chartData.frpRanges as range}
					<div class="flex items-center space-x-3">
						<span class="text-sm text-slate-600 w-24">{range.label}</span>
						<div class="flex-1 bg-slate-100 rounded-full h-6 relative">
							<div 
								class="bg-gradient-to-r from-blue-400 to-purple-500 h-6 rounded-full transition-all duration-300"
								style="width: {chartData.maxFRP > 0 ? (range.count / chartData.maxFRP) * 100 : 0}%"
							></div>
							<span class="absolute right-3 top-0 text-sm font-medium text-slate-700 leading-6">
								{range.count}
							</span>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>
{:else}
	<div class="text-center py-12">
		<BarChart3 class="w-12 h-12 text-slate-300 mx-auto mb-4" />
		<p class="text-slate-500">No data available for trend analysis</p>
	</div>
{/if}
