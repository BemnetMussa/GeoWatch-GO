<!-- src/routes/+page.svelte (Updated Version) -->

<script lang="ts">
	// 1. Change this import from MapContainer to our new GlobeView
	import GlobeView from '$lib/components/GlobalView.svelte';
	import MapContainer from '$lib/components/MapContainer.svelte';
	import AnalyticsDashboard from '$lib/components/AnalyticsDashboard.svelte';
	import NotificationSystem from '$lib/components/NotificationSystem.svelte';
	import type { FireData } from '$lib/types';
	import { onMount } from 'svelte';
	
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
<div class="h-screen w-screen flex flex-col bg-slate-900">
	<div class="flex-1 relative">
		{#if mounted}
			<!-- View switcher -->
			<div class="absolute top-4 right-4 z-20 flex items-center space-x-2 bg-slate-800 bg-opacity-90 text-white px-3 py-2 rounded-md">
				<button class="px-2 py-1 rounded {activeView === 'globe' ? 'bg-orange-600' : 'bg-transparent'}" onclick={() => activeView = 'globe'}>Globe</button>
				<button class="px-2 py-1 rounded {activeView === 'map' ? 'bg-orange-600' : 'bg-transparent'}" onclick={() => activeView = 'map'}>Map</button>
				<button class="px-2 py-1 rounded bg-slate-700" onclick={() => showAnalytics = true}>Analytics</button>
				<label class="flex items-center space-x-2 text-xs">
					<input type="checkbox" bind:checked={notificationsEnabled} />
					<span>Notifications</span>
				</label>
			</div>

			<!-- Main view -->
			<div class="absolute inset-0">
				{#if activeView === 'globe'}
					<GlobeView />
				{:else}
					<MapContainer />
				{/if}
			</div>

			<!-- Analytics modal (works for both views; MapContainer also has its own, this is demo for Globe) -->
			<AnalyticsDashboard fires={demoFires} isOpen={showAnalytics} onClose={() => showAnalytics = false} />

			<!-- Notification preview (optional) -->
			<NotificationSystem enabled={notificationsEnabled} newFires={demoNewFires} onNotificationClick={() => {}} />
		{:else}
			<!-- Loading -->
			<div class="h-full flex items-center justify-center">
				<div class="text-center">
					<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-600 mx-auto mb-4"></div>
					<p class="text-slate-400">Loading planet...</p>
				</div>
			</div>
		{/if}
	</div>
</div>