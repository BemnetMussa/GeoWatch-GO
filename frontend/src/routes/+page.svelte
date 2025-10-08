<!-- src/routes/+page.svelte (Updated Version) -->

<script lang="ts">
	// 1. Change this import from MapContainer to our new GlobeView
	import GlobeView from '$lib/components/GlobalView.svelte';
	import MapContainer from '$lib/components/MapContainer.svelte';
	import AnalyticsDashboard from '$lib/components/AnalyticsDashboard.svelte';
	import NotificationSystem from '$lib/components/NotificationSystem.svelte';
	import type { FireData } from '$lib/types';
	import { onMount } from 'svelte';
	import GlobalView from '$lib/components/GlobalView.svelte';
	
	let mounted = $state(false);
	let activeView: 'globe' | 'map' = $state('globe');
	let showAnalytics = $state(false);
	let notificationsEnabled = $state(false);
	let demoNewFires: FireData[] = $state([]);
	let demoFires: FireData[] = $state([
		{
			latitude: 34.05, longitude: -118.25, brightness: 320, scan: 1, track: 1,
			acq_date: '2025-10-05', acq_time: '1200', satellite: 'MODIS', confidence: 85,
			version: '1', bright_t31: 300, frp: 12.5, daynight: 'D'
		},
		{
			latitude: 37.77, longitude: -122.42, brightness: 295, scan: 1, track: 1,
			acq_date: '2025-10-05', acq_time: '1130', satellite: 'VIIRS', confidence: 62,
			version: '1', bright_t31: 285, frp: 8.1, daynight: 'D'
		}
	]);
	
	onMount(() => {
		mounted = true;
		// seed notifications demo (optional)
		demoNewFires = demoFires.slice(0, 1);
	});
</script>

<svelte:head>
	<title>Wildfire Detection & Change Analysis</title>
	<meta name="description" content="Real-time wildfire detection and geographical change analysis using satellite data" />
</svelte:head>

<!-- 
  Your existing TailwindCSS structure is great. The key is to make sure the
  container for the globe takes up the full space available. `h-full` and `w-full`
  on the child div will achieve this.
-->
<div class="flex flex-col bg-slate-900 h-[calc(100vh-4rem)] w-full">
	<div class="flex-1 relative h-full w-full">
		{#if mounted}
			<!-- Main view -->
			<div class="absolute inset-0 h-full w-full">
				<GlobalView />
				<!-- {#if activeView === 'globe'}
					<GlobeView />
				{:else}
					<MapContainer />
				{/if} -->
			</div>

			<AnalyticsDashboard fires={demoFires} isOpen={showAnalytics} onClose={() => showAnalytics = false} />
			<NotificationSystem enabled={notificationsEnabled} newFires={demoNewFires} onNotificationClick={() => {}} />
		{:else}
			<div class="h-full flex items-center justify-center">
				<div class="text-center">
					<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-600 mx-auto mb-4"></div>
					<p class="text-slate-400">Loading planet...</p>
				</div>
			</div>
		{/if}
	</div>
</div>
