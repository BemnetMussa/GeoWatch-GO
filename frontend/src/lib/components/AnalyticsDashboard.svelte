<script lang="ts">
	import { BarChart3, TrendingUp, MapPin, Calendar, Download, AlertTriangle, Flame, Activity } from 'lucide-svelte';
	import type { FireData } from '$lib/types';
	import FireChart from './FireChart.svelte';
	import GeographicAnalysis from './GeographicAnalysis.svelte';
	import RiskAssessment from './RiskAssessment.svelte';
	
	interface Props {
		fires: FireData[];
		isOpen: boolean;
		onClose: () => void;
	}
	
	let { fires, isOpen, onClose }: Props = $props();
	
	let activeTab = $state('overview');
	
	// Calculate analytics data
	let analytics = $derived.by(() => {
		if (!fires.length) return null;
		
		const totalFires = fires.length;
		const avgConfidence = fires.reduce((sum, fire) => sum + fire.confidence, 0) / totalFires;
		const avgBrightness = fires.reduce((sum, fire) => sum + fire.brightness, 0) / totalFires;
		const totalFRP = fires.reduce((sum, fire) => sum + (fire.frp || 0), 0);
		
		// Confidence distribution
		const highConfidence = fires.filter(f => f.confidence >= 80).length;
		const mediumConfidence = fires.filter(f => f.confidence >= 60 && f.confidence < 80).length;
		const lowConfidence = fires.filter(f => f.confidence < 60).length;
		
		// Satellite distribution
		const satelliteData = fires.reduce((acc, fire) => {
			acc[fire.satellite] = (acc[fire.satellite] || 0) + 1;
			return acc;
		}, {} as Record<string, number>);
		
		// Time distribution (by hour)
		const timeData = fires.reduce((acc, fire) => {
			const hour = fire.acq_time.substring(0, 2);
			acc[hour] = (acc[hour] || 0) + 1;
			return acc;
		}, {} as Record<string, number>);
		
		// Geographic hotspots
		const latBounds = { min: Math.min(...fires.map(f => f.latitude)), max: Math.max(...fires.map(f => f.latitude)) };
		const lngBounds = { min: Math.min(...fires.map(f => f.longitude)), max: Math.max(...fires.map(f => f.longitude)) };
		
		return {
			totalFires,
			avgConfidence: Math.round(avgConfidence * 10) / 10,
			avgBrightness: Math.round(avgBrightness * 10) / 10,
			totalFRP: Math.round(totalFRP * 10) / 10,
			confidenceDistribution: { high: highConfidence, medium: mediumConfidence, low: lowConfidence },
			satelliteData,
			timeData,
			geographicBounds: { lat: latBounds, lng: lngBounds }
		};
	});
	
	function exportData() {
		if (!fires.length) return;
		
		const csvContent = [
			'Latitude,Longitude,Brightness,Confidence,FRP,Satellite,Date,Time',
			...fires.map(fire => 
				`${fire.latitude},${fire.longitude},${fire.brightness},${fire.confidence},${fire.frp || 0},${fire.satellite},${fire.acq_date},${fire.acq_time}`
			)
		].join('\n');
		
		const blob = new Blob([csvContent], { type: 'text/csv' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `fire-data-${new Date().toISOString().split('T')[0]}.csv`;
		a.click();
		URL.revokeObjectURL(url);
	}
</script>

{#if isOpen}
    <!-- Centered modal with subtle transparent backdrop, minimal obstruction -->
    <div class="fixed inset-0 z-40 pointer-events-none">
        <div class="absolute inset-0 bg-black/10"></div>
    </div>
    <div class="fixed inset-0 z-50 flex items-start justify-center p-6">
        <div class="bg-white/95 backdrop-blur rounded-lg shadow-xl w-full max-w-6xl max-h-[90vh] overflow-hidden pointer-events-auto">
			<!-- Header -->
			<div class="flex items-center justify-between p-6 border-b border-slate-200">
				<div class="flex items-center space-x-3">
					<BarChart3 class="w-6 h-6 text-orange-600" />
					<h2 class="text-xl font-bold text-slate-900">Fire Analytics Dashboard</h2>
				</div>
				<div class="flex items-center space-x-2">
					<button
						onclick={exportData}
						class="flex items-center space-x-2 px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 rounded-md transition-colors"
					>
						<Download class="w-4 h-4" />
						<span>Export Data</span>
					</button>
					<button
                        onclick={onClose}
                        aria-label="Close"
                        class="text-slate-400 hover:text-slate-600 transition-colors p-2"
                        >
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                    </button>

				</div>
			</div>
			
			<!-- Tabs -->
			<div class="border-b border-slate-200">
				<nav class="flex space-x-8 px-6">
					<button
						onclick={() => activeTab = 'overview'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'overview' ? 'border-orange-500 text-orange-600' : 'border-transparent text-slate-500 hover:text-slate-700'}"
					>
						Overview
					</button>
					<button
						onclick={() => activeTab = 'trends'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'trends' ? 'border-orange-500 text-orange-600' : 'border-transparent text-slate-500 hover:text-slate-700'}"
					>
						Trends
					</button>
					<button
						onclick={() => activeTab = 'geographic'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'geographic' ? 'border-orange-500 text-orange-600' : 'border-transparent text-slate-500 hover:text-slate-700'}"
					>
						Geographic
					</button>
					<button
						onclick={() => activeTab = 'risk'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab === 'risk' ? 'border-orange-500 text-orange-600' : 'border-transparent text-slate-500 hover:text-slate-700'}"
					>
						Risk Assessment
					</button>
				</nav>
			</div>
			
			<!-- Content -->
			<div class="p-6 overflow-y-auto max-h-[calc(90vh-200px)]">
				{#if !analytics}
					<div class="text-center py-12">
						<Flame class="w-12 h-12 text-slate-300 mx-auto mb-4" />
						<p class="text-slate-500">No fire data available for analysis</p>
					</div>
				{:else if activeTab === 'overview'}
					<!-- Overview Tab -->
					<div class="space-y-6">
						<!-- Key Metrics -->
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
							<div class="bg-gradient-to-br from-red-50 to-red-100 rounded-lg p-6">
								<div class="flex items-center justify-between">
									<div>
										<p class="text-sm font-medium text-red-600">Total Fires</p>
										<p class="text-2xl font-bold text-red-900">{analytics.totalFires}</p>
									</div>
									<Flame class="w-8 h-8 text-red-500" />
								</div>
							</div>
							
							<div class="bg-gradient-to-br from-orange-50 to-orange-100 rounded-lg p-6">
								<div class="flex items-center justify-between">
									<div>
										<p class="text-sm font-medium text-orange-600">Avg Confidence</p>
										<p class="text-2xl font-bold text-orange-900">{analytics.avgConfidence}%</p>
									</div>
									<TrendingUp class="w-8 h-8 text-orange-500" />
								</div>
							</div>
							
							<div class="bg-gradient-to-br from-yellow-50 to-yellow-100 rounded-lg p-6">
								<div class="flex items-center justify-between">
									<div>
										<p class="text-sm font-medium text-yellow-600">Avg Brightness</p>
										<p class="text-2xl font-bold text-yellow-900">{analytics.avgBrightness}K</p>
									</div>
									<Activity class="w-8 h-8 text-yellow-500" />
								</div>
							</div>
							
							<div class="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-6">
								<div class="flex items-center justify-between">
									<div>
										<p class="text-sm font-medium text-blue-600">Total FRP</p>
										<p class="text-2xl font-bold text-blue-900">{analytics.totalFRP} MW</p>
									</div>
									<BarChart3 class="w-8 h-8 text-blue-500" />
								</div>
							</div>
						</div>
						
						<!-- Confidence Distribution -->
						<div class="bg-white rounded-lg border border-slate-200 p-6">
							<h3 class="text-lg font-semibold text-slate-900 mb-4">Confidence Distribution</h3>
							<div class="space-y-4">
								<div class="flex items-center justify-between">
									<span class="text-sm text-slate-600">High Confidence (≥80%)</span>
									<div class="flex items-center space-x-2">
										<div class="w-32 bg-slate-200 rounded-full h-2">
											<div class="bg-red-600 h-2 rounded-full" style="width: {(analytics.confidenceDistribution.high / analytics.totalFires) * 100}%"></div>
										</div>
										<span class="text-sm font-medium text-slate-900">{analytics.confidenceDistribution.high}</span>
									</div>
								</div>
								<div class="flex items-center justify-between">
									<span class="text-sm text-slate-600">Medium Confidence (60-79%)</span>
									<div class="flex items-center space-x-2">
										<div class="w-32 bg-slate-200 rounded-full h-2">
											<div class="bg-orange-500 h-2 rounded-full" style="width: {(analytics.confidenceDistribution.medium / analytics.totalFires) * 100}%"></div>
										</div>
										<span class="text-sm font-medium text-slate-900">{analytics.confidenceDistribution.medium}</span>
									</div>
								</div>
								<div class="flex items-center justify-between">
				                    <span class="text-sm text-slate-600">Low Confidence (&lt;60%)</span> 

									<div class="flex items-center space-x-2">
										<div class="w-32 bg-slate-200 rounded-full h-2">
											<div class="bg-yellow-500 h-2 rounded-full" style="width: {(analytics.confidenceDistribution.low / analytics.totalFires) * 100}%"></div>
										</div>
										<span class="text-sm font-medium text-slate-900">{analytics.confidenceDistribution.low}</span>
									</div>
								</div>
							</div>
						</div>
						
						<!-- Satellite Sources -->
						<div class="bg-white rounded-lg border border-slate-200 p-6">
							<h3 class="text-lg font-semibold text-slate-900 mb-4">Data Sources</h3>
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
								{#each Object.entries(analytics.satelliteData) as [satellite, count]}
									<div class="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
										<span class="text-sm font-medium text-slate-700">{satellite}</span>
										<span class="text-sm text-slate-900 bg-white px-2 py-1 rounded">{count}</span>
									</div>
								{/each}
							</div>
						</div>
					</div>
				{:else if activeTab === 'trends'}
					<FireChart {fires} />
				{:else if activeTab === 'geographic'}
					<GeographicAnalysis {fires} />
				{:else if activeTab === 'risk'}
					<RiskAssessment {fires} />
				{/if}
			</div>
		</div>
	</div>
{/if}
