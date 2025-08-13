<script lang="ts">
	import { Flame, Satellite, BarChart3, Settings } from 'lucide-svelte';
</script>

<header class="bg-white border-b border-slate-200 shadow-sm">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex justify-between items-center h-16">
			<div class="flex items-center space-x-3">
				<div class="flex items-center justify-center w-10 h-10 bg-gradient-to-br rounded-lg">
					  <Satellite className="h-6 w-6 text-primary" />
				</div>
				<div>
					<h1 class="text-xl font-bold text-slate-900">GeoWatch</h1>
					<p class="text-sm text-slate-600">Geospatial Change Detection</p>
				</div>
			</div>
			
			<nav class="hidden md:flex items-center space-x-6">
				<button class="flex items-center space-x-2 px-3 py-2 text-sm font-medium text-slate-700 hover:text-orange-600 transition-colors">
					<Satellite class="w-4 h-4" />
					<span>Live Data</span>
				</button>
				<button class="flex items-center space-x-2 px-3 py-2 text-sm font-medium text-slate-700 hover:text-orange-600 transition-colors">
					<BarChart3 class="w-4 h-4" />
					<span>Analytics</span>
				</button>
				<button class="flex items-center space-x-2 px-3 py-2 text-sm font-medium text-slate-700 hover:text-orange-600 transition-colors">
					<Settings class="w-4 h-4" />
					<span>Settings</span>
				</button>
			</nav>
		</div>
	</div>
</header>
