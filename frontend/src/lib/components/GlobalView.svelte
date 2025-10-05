<!-- src/lib/components/GlobeView.svelte -->

<script lang="ts">
	import { onMount } from 'svelte';
	import { env as publicEnv } from '$env/dynamic/public';
	const BACKEND_URL = publicEnv.PUBLIC_GO_BACKEND_URL ?? 'http://localhost:8081';
	
	// --- 1. THE CORRECTED IMPORT BLOCK ---
	import { 
		Viewer, 
		Ion, 
		createWorldTerrainAsync,
        UrlTemplateImageryProvider,
		Rectangle,
		Entity,
		Color,
		Math as CesiumMath, // Aliased to avoid conflicts
		ClassificationType,   // <-- IMPORTED
		ImageMaterialProperty, // <-- IMPORTED
		ImageryLayer,
		Credit,
		SingleTileImageryProvider,
		buildModuleUrl
	} from 'cesium';
    import MapControls from './MapControls.svelte';
    import ChangeDetectionPanel from './ChangeDetectionPanel.svelte';
    import RealTimeSettings from './RealTimeSettings.svelte';
    import NotificationSystem from './NotificationSystem.svelte';
    import AnalyticsDashboard from './AnalyticsDashboard.svelte';
    import type { FireData } from '$lib/types';
	import type { Viewer as CesiumViewer } from 'cesium';
	import 'cesium/Build/Cesium/Widgets/widgets.css';

	let cesiumContainer: HTMLElement;
	let viewer: CesiumViewer | null = null;
    let loadingAnalysis = false;
    let startYear: number = 2018;
    let endYear: number = 2023;
	
	// This will store our analysis layer so we can remove it later
	let changeLayerEntity: Entity | null = null;
	let changeImageryLayer: ImageryLayer | null = null; 
	let lastImageUrl: string | null = null;
	let overlayOpacity: number = 0.75;

	// MapControls-like state for Globe
	let activeLayer: 'satellite' | 'terrain' | 'street' = 'satellite';
	let showFires = true; // Placeholder (no fire layer on globe yet)
	let showChanges = false; // Placeholder (points change markers not on globe yet)
	let loading = false;
	let fireCount = 0; // Placeholder

	// Real-time/notifications/analytics state
	let autoRefresh = false;
	let refreshIntervalMs = 300000;
	let notificationsEnabled = false;
	let connectionStatus: 'connected' | 'disconnected' | 'syncing' = 'disconnected';
	let lastUpdate: Date | null = null;
	let updateCount = 0;
	let showAnalytics = false;
	let demoFires: FireData[] = [
		{
			latitude: 34.05, longitude: -118.25, brightness: 320, scan: 1, track: 1,
			acq_date: '2025-10-05', acq_time: '1200', satellite: 'MODIS', confidence: 85,
			version: '1', bright_t31: 300, frp: 12.5, daynight: 'D'
		}
	];
	let demoNewFires: FireData[] = demoFires;

	// Keep references to base layer and labels overlay
	let baseImageryLayer: ImageryLayer | null = null;
	let labelsImageryLayer: ImageryLayer | null = null;

	// Toast notification state
	let toast: { type: 'info' | 'success' | 'error'; message: string } | null = null;
	let toastTimer: number | null = null;

	function notify(message: string, type: 'info' | 'success' | 'error' = 'info', durationMs: number = 4000) {
		toast = { type, message };
		if (toastTimer) {
			clearTimeout(toastTimer);
			toastTimer = null;
		}
		toastTimer = setTimeout(() => {
			toast = null;
			toastTimer = null;
		}, durationMs) as unknown as number;
	}

    // --- NEW: Zoom safeguard variables ---
	let isZoomedInEnough = false;
	let currentAltitude = 0;
	let viewportArea = 0;
	const MAX_ANALYSIS_ALTITUDE = 1000000; // 1000km - adjust this threshold as needed
	const MAX_ANALYSIS_AREA = 100000; // 100,000 km² - adjust this threshold as needed

	// --- NEW: Function to check if current zoom level is suitable for analysis ---
	function checkZoomLevel() {
		if (!viewer) return;
		
		const camera = viewer.camera;
		const position = camera.position;
		const ellipsoid = viewer.scene.globe.ellipsoid;
		
		// Get the camera's height above the ellipsoid
		const cartographic = ellipsoid.cartesianToCartographic(position);
		currentAltitude = cartographic.height;
		
		// Check if we're zoomed in enough for analysis
		isZoomedInEnough = currentAltitude <= MAX_ANALYSIS_ALTITUDE;
		
		// Calculate current viewport area
		viewportArea = calculateViewportArea();
		
		console.log(`Camera altitude: ${(currentAltitude / 1000).toFixed(1)}km, Analysis enabled: ${isZoomedInEnough}, Viewport area: ${viewportArea.toFixed(1)} km²`);
	}

	// --- NEW: Reactive statement to update viewport area when needed ---
	$: if (viewer && isZoomedInEnough !== undefined) {
		viewportArea = calculateViewportArea();
	}

	// --- NEW: Function to calculate the area of the current viewport ---
	function calculateViewportArea(): number {
		if (!viewer) return 0;
		
		const viewRect = viewer.camera.computeViewRectangle();
		if (!viewRect) return 0;
		
		const west = CesiumMath.toDegrees(viewRect.west);
		const south = CesiumMath.toDegrees(viewRect.south);
		const east = CesiumMath.toDegrees(viewRect.east);
		const north = CesiumMath.toDegrees(viewRect.north);
		
		// Calculate approximate area in square kilometers
		const latDiff = Math.abs(north - south);
		const lonDiff = Math.abs(east - west);
		const avgLat = (north + south) / 2;
		
		// Convert to kilometers (rough approximation)
		const latKm = latDiff * 111.32; // 1 degree latitude ≈ 111.32 km
		const lonKm = lonDiff * (111.32 * Math.cos(avgLat * Math.PI / 180));
		
		return latKm * lonKm;
	}

	// The initializeGlobe function is already correct from before
	async function initializeGlobe() {
		if (!cesiumContainer) return;
		// set base URL for Cesium static assets
		try {
			// Some Cesium builds expose setBaseUrl; guard in case it's absent
			(buildModuleUrl as any).setBaseUrl?.('/Cesium/');
		} catch {}

    // Cesium Ion token from env
    Ion.defaultAccessToken = publicEnv.PUBLIC_CESIUM_ION_TOKEN ?? '';
		viewer = new Viewer(cesiumContainer, {
			terrainProvider: await createWorldTerrainAsync(),
			animation: false, timeline: false, geocoder: false, homeButton: false, sceneModePicker: false,
			baseLayerPicker: false, navigationHelpButton: false, infoBox: true, selectionIndicator: false,
		});
		const labelCredit = new Credit('Map tiles by Stamen Design, under CC BY 3.0. Data by OpenStreetMap.', true);
		const labelImageryProvider = new UrlTemplateImageryProvider({
			url: 'https://tiles.stadiamaps.com/tiles/stamen_toner_labels/{z}/{x}/{y}.png',
			credit: labelCredit 
		});
		labelsImageryLayer = viewer.imageryLayers.addImageryProvider(labelImageryProvider);

		// Initialize base imagery as Satellite (Esri World Imagery)
		switchBaseLayer('satellite');
		viewer.scene.globe.enableLighting = true;
		viewer.scene.fog.enabled = true;
		viewer.scene.globe.depthTestAgainstTerrain = true;
		
		// --- NEW: Add camera movement listener for zoom checking ---
		viewer.camera.moveEnd.addEventListener(() => {
			checkZoomLevel();
		});
		
		// --- NEW: Initial zoom check ---
		checkZoomLevel();
	}

	function switchBaseLayer(layer: 'satellite' | 'terrain' | 'street') {
		if (!viewer) return;
		// Remove existing base (keep labels and change overlay)
		if (baseImageryLayer) {
			viewer.imageryLayers.remove(baseImageryLayer);
			baseImageryLayer = null;
		}
		let provider: any;
		if (layer === 'satellite') {
			provider = new UrlTemplateImageryProvider({
				url: 'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
				credit: new Credit('Esri World Imagery')
			});
		} else if (layer === 'terrain') {
			provider = new UrlTemplateImageryProvider({
				url: 'https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png',
				credit: new Credit('OpenTopoMap')
			});
		} else {
			provider = new UrlTemplateImageryProvider({
				url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
				credit: new Credit('OpenStreetMap')
			});
		}
		baseImageryLayer = viewer.imageryLayers.addImageryProvider(provider, 0);
		activeLayer = layer;
	}

	function switchLayer(layerName: 'satellite' | 'terrain' | 'street') {
		switchBaseLayer(layerName);
	}

	function toggleFires() {
		showFires = !showFires;
		// Placeholder: Globe fire overlays not implemented
	}

	function toggleChanges() {
		showChanges = !showChanges;
		// Placeholder: point change markers not implemented on globe
	}

	async function refreshData() {
		loading = true;
		try {
			// Placeholder refresh; mark connected
			connectionStatus = 'connected';
			lastUpdate = new Date();
			updateCount++;
		} finally {
			loading = false;
		}
	}

	function toggleAutoRefresh() {
		autoRefresh = !autoRefresh;
		// Placeholder: no interval loop here to avoid timers
	}

	function changeRefreshInterval(interval: number) {
		refreshIntervalMs = interval;
	}

	function toggleNotifications() {
		notificationsEnabled = !notificationsEnabled;
	}

	// --- 2. THE CORRECTED handleAnalyze FUNCTION ---
	async function handleAnalyze() {
		if (!viewer) return;
		
		// --- NEW: Zoom safeguard check ---
		if (!isZoomedInEnough) {
			const area = calculateViewportArea();
			notify(`Please zoom in further. Area ${area.toFixed(1)} km² is too large at current altitude.`, 'error');
			return;
		}
		
		// --- NEW: Area safeguard check ---
		const area = calculateViewportArea();
		if (area > MAX_ANALYSIS_AREA) {
			notify(`Viewport area ${area.toFixed(1)} km² exceeds recommended ${MAX_ANALYSIS_AREA.toFixed(0)} km². Please zoom in.`, 'error');
			return;
		}
		
		loadingAnalysis = true;

		 if (changeImageryLayer) {
        viewer.imageryLayers.remove(changeImageryLayer);
        changeImageryLayer = null;
			}

		// --- FIX FOR ERRORS 1 & 2 ---
		// Correctly remove the previous entity if it exists
		if (changeLayerEntity) {
			viewer.entities.remove(changeLayerEntity);
			changeLayerEntity = null;
		}

		try {
			const viewRect = viewer.camera.computeViewRectangle();
			if (!viewRect) {
				notify("Could not determine map bounds. Please zoom in closer.", 'error');
				loadingAnalysis = false;
				return;
			}
			
			// --- NEW: Calculate and log the analysis area ---
			const area = calculateViewportArea();
			console.log(`Starting analysis for viewport area: ${area.toFixed(1)} km²`);
			
			const west = CesiumMath.toDegrees(viewRect.west);
			const south = CesiumMath.toDegrees(viewRect.south);
			const east = CesiumMath.toDegrees(viewRect.east);
			const north = CesiumMath.toDegrees(viewRect.north);
			
			const aoi = [[ [west, south], [east, south], [east, north], [west, north], [west, south] ]];

			const response = await fetch(`${BACKEND_URL}/api/v1/changes`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ aoi: aoi, startDate: `${startYear}-01-01`, endDate: `${endYear}-12-31` })
			});

			if (!response.ok) {
				const errorText = await response.text();
				throw new Error(`Backend API call failed: ${errorText}`);
			}

			const imageBlob = await response.blob();
			// 1. Create an Image element from the blob URL
			const imageUrl = URL.createObjectURL(imageBlob);
			const image = new Image();
			lastImageUrl = imageUrl;
						
			// --- THIS IS THE NEW, CORRECT VISUALIZATION LOGIC ---
			
			// 1. Create a SingleTileImageryProvider. This object knows how to handle a single image tile.
			        // Our Python backend is configured to always return a 2048x2048 image.
			const imageWidth = 2048;
			const imageHeight = 2048;

			const imageryProvider = new SingleTileImageryProvider({
				url: imageUrl,
				rectangle: viewRect,
				
				// --- THIS IS THE FIX ---
				// We must explicitly tell the provider the dimensions of our image.
				tileWidth: imageWidth,
				tileHeight: imageHeight
			});

			// 2. Add this provider to the globe as a new ImageryLayer.
			// We add it at a specific index (1) to make sure it appears on top of the base map.
			changeImageryLayer = viewer.imageryLayers.addImageryProvider(imageryProvider);
			// ensure analysis overlay is on top of all imagery (including labels)
			viewer.imageryLayers.raiseToTop(changeImageryLayer);
			
			// 3. Set the opacity of the entire layer.
			changeImageryLayer.alpha = overlayOpacity;

			// 4. Fly to the new layer's extent.
			viewer.flyTo(changeImageryLayer);

			console.log("Successfully added GEE result as an ImageryLayer.");
			notify('Analysis complete. Overlay added.', 'success');

		} catch (error) {
			console.error("Analysis failed:", error);
			notify(`Analysis failed. See console for details.`, 'error');
		} finally {
			loadingAnalysis = false;
		}
	}
	onMount(() => {
		initializeGlobe();
		return () => {
			if (viewer) {
				viewer.destroy();
			}
		};
	});
</script>

<!-- The HTML and CSS are perfect. No changes needed. -->

<div class="relative w-full h-full">
	<div bind:this={cesiumContainer} class="w-full h-full"></div>
    <MapControls 
		activeLayer={activeLayer}
		showFires={showFires}
		loading={loading}
		fireCount={fireCount}
		showChanges={showChanges}
		changeCount={0}
		onLayerChange={switchLayer}
		onToggleFires={toggleFires}
		onToggleChanges={toggleChanges}
		onRefresh={refreshData}
		onShowAnalytics={() => showAnalytics = true}
	/>

<!-- Year selection removed; years are now in ChangeDetectionPanel -->

<!-- Overlay Controls: opacity + download (single card) -->
{#if changeImageryLayer}
	<div class="absolute bottom-4 left-4 z-10 bg-slate-800 bg-opacity-90 text-white rounded-lg shadow-lg p-3 w-72 space-y-3">
		<div class="text-xs">Overlay Opacity: {(overlayOpacity * 100).toFixed(0)}%</div>
		<input
			type="range"
			min="0"
			max="1"
			step="0.05"
			bind:value={overlayOpacity}
			on:input={() => { if (changeImageryLayer) changeImageryLayer.alpha = overlayOpacity; }}
			class="w-full"
		/>
		{#if lastImageUrl}
			<a 
				href={lastImageUrl} 
				download="analysis_result.png"
				class="inline-flex items-center justify-center bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-3 rounded-md text-xs"
			>
				💾 Download Analysis
			</a>
		{/if}
	</div>
{/if}

<!-- Change Detection panel (uses Globe handler with days-to-years mapping) -->
<ChangeDetectionPanel 
    onAnalyze={async () => { await handleAnalyze(); }}
    result={null}
    loading={loadingAnalysis}
    yearsMode={true}
    startYear={startYear}
    endYear={endYear}
    onYearsChange={(s, e) => { startYear = s; endYear = e; }}
/>

<!-- Real-time settings -->
<RealTimeSettings 
	autoRefresh={autoRefresh}
	refreshInterval={refreshIntervalMs}
	notifications={notificationsEnabled}
	connectionStatus={connectionStatus}
	lastUpdate={lastUpdate}
	updateCount={updateCount}
	onToggleAutoRefresh={toggleAutoRefresh}
	onChangeInterval={changeRefreshInterval}
	onToggleNotifications={toggleNotifications}
/>

<!-- Notifications and analytics -->
<NotificationSystem enabled={notificationsEnabled} newFires={demoNewFires} onNotificationClick={() => {}} />
<AnalyticsDashboard fires={demoFires} isOpen={showAnalytics} onClose={() => showAnalytics = false} />

<!-- Optimal area/zoom helper popup (center top) -->
<div class="absolute top-4 left-1/2 -translate-x-1/2 z-10 p-3 bg-slate-800 bg-opacity-90 rounded-lg shadow-lg text-center max-w-sm">
    <div class="text-sm text-slate-300 mb-2">
        {isZoomedInEnough ? '🎯 Ready for Analysis' : '🔍 Adjusting Viewport'}
    </div>
	<div class="text-xs text-slate-400 space-y-1">
        {#if isZoomedInEnough}
			<div>✅ Zoom level optimal</div>
			<div>✅ Viewport area: {viewportArea.toFixed(1)} km²</div>
        {:else}
			<div>⚠️ Current altitude: {(currentAltitude / 1000).toFixed(1)} km</div>
			<div>⚠️ Target: &lt; {(MAX_ANALYSIS_ALTITUDE / 1000).toFixed(1)} km</div>
			<div>💡 Zoom in until you can clearly see features</div>
        {/if}
    </div>
</div>

<!-- Toast -->
{#if toast}
	<div class="absolute bottom-4 left-1/2 -translate-x-1/2 z-20">
		<div class="px-4 py-3 rounded-lg shadow-lg text-sm font-medium
			{toast.type === 'success' ? 'bg-green-600 text-white' : ''}
			{toast.type === 'error' ? 'bg-red-600 text-white' : ''}
			{toast.type === 'info' ? 'bg-slate-800 text-white' : ''}">
			{toast.message}
		</div>
	</div>
{/if}
</div>
<style>
    :global(.cesium-widget canvas) {
            width: 100%;
            height: 100%;
        }
</style>
