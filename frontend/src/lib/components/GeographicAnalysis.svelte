<script lang="ts">
	import { MapPin, Globe, Target } from 'lucide-svelte';
	import type { FireData } from '$lib/types';
	
	interface Props {
		fires: FireData[];
	}
	
	let { fires }: Props = $props();
	
	let geoAnalysis = $derived.by(() => {
		if (!fires.length) return null;
		
		// Calculate geographic bounds
		const latitudes = fires.map(f => f.latitude);
		const longitudes = fires.map(f => f.longitude);
		
		const bounds = {
			north: Math.max(...latitudes),
			south: Math.min(...latitudes),
			east: Math.max(...longitudes),
			west: Math.min(...longitudes),
			center: {
				lat: (Math.max(...latitudes) + Math.min(...latitudes)) / 2,
				lng: (Math.max(...longitudes) + Math.min(...longitudes)) / 2
			}
		};
		
		// Create density grid (simplified)
		const gridSize = 0.5; // degrees
		const densityGrid: Record<string, { count: number; avgBrightness: number; fires: FireData[] }> = {};
		
		fires.forEach(fire => {
			const gridLat = Math.floor(fire.latitude / gridSize) * gridSize;
			const gridLng = Math.floor(fire.longitude / gridSize) * gridSize;
			const key = `${gridLat},${gridLng}`;
			
			if (!densityGrid[key]) {
				densityGrid[key] = { count: 0, avgBrightness: 0, fires: [] };
			}
			
			densityGrid[key].count++;
			densityGrid[key].fires.push(fire);
		});
		
		// Calculate average brightness for each grid cell
		Object.values(densityGrid).forEach(cell => {
			cell.avgBrightness = cell.fires.reduce((sum, fire) => sum + fire.brightness, 0) / cell.count;
		});
		
		// Find hotspots (top 10 grid cells by fire count)
		const hotspots = Object.entries(densityGrid)
			.map(([key, data]) => {
				const [lat, lng] = key.split(',').map(Number);
				return {
					lat,
					lng,
					count: data.count,
					avgBrightness: Math.round(data.avgBrightness * 10) / 10,
					fires: data.fires
				};
			})
			.sort((a, b) => b.count - a.count)
			.slice(0, 10);
		
		// State/region analysis (simplified - based on lat/lng ranges)
		const regions = [
			{ name: 'Western US', bounds: { latMin: 32, latMax: 49, lngMin: -125, lngMax: -100 }, count: 0 },
			{ name: 'Central US', bounds: { latMin: 25, latMax: 49, lngMin: -100, lngMax: -90 }, count: 0 },
			{ name: 'Eastern US', bounds: { latMin: 25, latMax: 49, lngMin: -90, lngMax: -65 }, count: 0 },
			{ name: 'Alaska', bounds: { latMin: 55, latMax: 72, lngMin: -180, lngMax: -130 }, count: 0 },
			{ name: 'Other', bounds: { latMin: -90, latMax: 90, lngMin: -180, lngMax: 180 }, count: 0 }
		];
		
		fires.forEach(fire => {
			let assigned = false;
			for (const region of regions.slice(0, -1)) { // Exclude 'Other' from initial check
				if (fire.latitude >= region.bounds.latMin && fire.latitude <= region.bounds.latMax &&
					fire.longitude >= region.bounds.lngMin && fire.longitude <= region.bounds.lngMax) {
					region.count++;
					assigned = true;
					break;
				}
			}
			if (!assigned) {
				regions[regions.length - 1].count++; // Add to 'Other'
			}
		});
		
		return {
			bounds,
			hotspots,
			regions: regions.filter(r => r.count > 0),
			totalArea: Math.abs(bounds.north - bounds.south) * Math.abs(bounds.east - bounds.west)
		};
	});
</script>

{#if geoAnalysis}
	<div class="space-y-8">
		<!-- Geographic Summary -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-green-600">Coverage Area</p>
						<p class="text-2xl font-bold text-green-900">{geoAnalysis.totalArea.toFixed(1)}°²</p>
					</div>
					<Globe class="w-8 h-8 text-green-500" />
				</div>
			</div>
			
			<div class="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-purple-600">Hotspots</p>
						<p class="text-2xl font-bold text-purple-900">{geoAnalysis.hotspots.length}</p>
					</div>
					<Target class="w-8 h-8 text-purple-500" />
				</div>
			</div>
			
			<div class="bg-gradient-to-br from-indigo-50 to-indigo-100 rounded-lg p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-indigo-600">Center Point</p>
						<p class="text-sm font-bold text-indigo-900">
							{geoAnalysis.bounds.center.lat.toFixed(2)}°, {geoAnalysis.bounds.center.lng.toFixed(2)}°
						</p>
					</div>
					<MapPin class="w-8 h-8 text-indigo-500" />
				</div>
			</div>
		</div>
		
		<!-- Regional Distribution -->
		<div class="bg-white rounded-lg border border-slate-200 p-6">
			<h3 class="text-lg font-semibold text-slate-900 mb-6">Regional Distribution</h3>
			<div class="space-y-4">
				{#each geoAnalysis.regions as region}
					<div class="flex items-center justify-between">
						<span class="text-sm font-medium text-slate-700">{region.name}</span>
						<div class="flex items-center space-x-3">
							<div class="w-32 bg-slate-200 rounded-full h-3">
								<div 
									class="bg-gradient-to-r from-indigo-400 to-purple-500 h-3 rounded-full"
									style="width: {(region.count / fires.length) * 100}%"
								></div>
							</div>
							<span class="text-sm font-bold text-slate-900 w-8">{region.count}</span>
						</div>
					</div>
				{/each}
			</div>
		</div>
		
		<!-- Top Hotspots -->
		<div class="bg-white rounded-lg border border-slate-200 p-6">
			<h3 class="text-lg font-semibold text-slate-900 mb-6">Fire Density Hotspots</h3>
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b border-slate-200">
							<th class="text-left py-3 px-2 font-medium text-slate-600">Rank</th>
							<th class="text-left py-3 px-2 font-medium text-slate-600">Location</th>
							<th class="text-left py-3 px-2 font-medium text-slate-600">Fire Count</th>
							<th class="text-left py-3 px-2 font-medium text-slate-600">Avg Brightness</th>
							<th class="text-left py-3 px-2 font-medium text-slate-600">Density</th>
						</tr>
					</thead>
					<tbody>
						{#each geoAnalysis.hotspots as hotspot, index}
							<tr class="border-b border-slate-100 hover:bg-slate-50">
								<td class="py-3 px-2">
									<span class="inline-flex items-center justify-center w-6 h-6 bg-orange-100 text-orange-800 text-xs font-bold rounded-full">
										{index + 1}
									</span>
								</td>
								<td class="py-3 px-2 font-mono text-xs">
									{hotspot.lat.toFixed(2)}°, {hotspot.lng.toFixed(2)}°
								</td>
								<td class="py-3 px-2 font-semibold">{hotspot.count}</td>
								<td class="py-3 px-2">{hotspot.avgBrightness}K</td>
								<td class="py-3 px-2">
									<div class="flex items-center space-x-2">
										<div class="w-16 bg-slate-200 rounded-full h-2">
											<div 
												class="bg-red-500 h-2 rounded-full"
												style="width: {(hotspot.count / geoAnalysis.hotspots[0].count) * 100}%"
											></div>
										</div>
										<span class="text-xs text-slate-500">
											{((hotspot.count / geoAnalysis.hotspots[0].count) * 100).toFixed(0)}%
										</span>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
{:else}
	<div class="text-center py-12">
		<MapPin class="w-12 h-12 text-slate-300 mx-auto mb-4" />
		<p class="text-slate-500">No geographic data available for analysis</p>
	</div>
{/if}
