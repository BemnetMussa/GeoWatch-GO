<script lang="ts">
	import { X, Flame, Thermometer, Zap, Satellite } from 'lucide-svelte';
	import type { FireData } from '$lib/types';
	
	interface Props {
		fire: FireData;
		onClose: () => void;
	}
	
	let { fire, onClose }: Props = $props();
	
	function getConfidenceColor(confidence: number) {
		if (confidence >= 80) return 'text-red-600 bg-red-50';
		if (confidence >= 60) return 'text-orange-600 bg-orange-50';
		if (confidence >= 40) return 'text-yellow-600 bg-yellow-50';
		return 'text-slate-600 bg-slate-50';
	}
	
	function getConfidenceLabel(confidence: number) {
		if (confidence >= 80) return 'Very High';
		if (confidence >= 60) return 'High';
		if (confidence >= 40) return 'Medium';
		return 'Low';
	}
</script>

<div class="absolute bottom-4 left-4 z-10 bg-white rounded-lg shadow-lg border border-slate-200 p-4 max-w-sm">
	<div class="flex items-start justify-between mb-3">
		<div class="flex items-center space-x-2">
			<Flame class="w-5 h-5 text-orange-600" />
			<h3 class="font-semibold text-slate-900">Fire Detection</h3>
		</div>
		<button
			onclick={onClose}
			class="text-slate-400 hover:text-slate-600 transition-colors"
		>
			<X class="w-4 h-4" />
		</button>
	</div>
	
	<div class="space-y-3">
		<!-- Confidence -->
		<div class="flex items-center justify-between">
			<span class="text-sm text-slate-600">Confidence</span>
			<div class="flex items-center space-x-2">
				<span class="text-sm font-medium {getConfidenceColor(fire.confidence)} px-2 py-1 rounded">
					{getConfidenceLabel(fire.confidence)}
				</span>
				<span class="text-sm text-slate-500">{fire.confidence}%</span>
			</div>
		</div>
		
		<!-- Brightness -->
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-1">
				<Thermometer class="w-4 h-4 text-slate-400" />
				<span class="text-sm text-slate-600">Brightness</span>
			</div>
			<span class="text-sm font-medium text-slate-900">{fire.brightness.toFixed(1)}K</span>
		</div>
		
		<!-- Fire Radiative Power -->
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-1">
				<Zap class="w-4 h-4 text-slate-400" />
				<span class="text-sm text-slate-600">FRP</span>
			</div>
			<span class="text-sm font-medium text-slate-900">{(fire.frp || 0).toFixed(1)} MW</span>
		</div>
		
		<!-- Satellite -->
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-1">
				<Satellite class="w-4 h-4 text-slate-400" />
				<span class="text-sm text-slate-600">Satellite</span>
			</div>
			<span class="text-sm font-medium text-slate-900">{fire.satellite}</span>
		</div>
		
		<!-- Location -->
		<div class="pt-2 border-t border-slate-100">
			<div class="text-xs text-slate-500 space-y-1">
				<div>Lat: {fire.latitude.toFixed(4)}°</div>
				<div>Lng: {fire.longitude.toFixed(4)}°</div>
				<div>Date: {fire.acq_date} {fire.acq_time}</div>
			</div>
		</div>
	</div>
</div>
