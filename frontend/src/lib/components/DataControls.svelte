<!-- frontend/src/lib/components/DateControls.svelte -->

<script lang="ts">
	// These are the two-way bindings for the year values
	export let startYear: number;
	export let endYear: number;
</script>

<div class="absolute bottom-4 left-1/2 -translate-x-1/2 z-10 p-3 bg-slate-800 bg-opacity-80 rounded-lg shadow-lg flex items-center gap-4">
	<!-- "Before" Year Input -->
	<div class="flex flex-col items-center">
		<label for="start-year" class="text-xs font-bold text-slate-300 mb-1">BEFORE</label>
		<input 
			id="start-year"
			type="number" 
			bind:value={startYear}
			class="bg-slate-900 text-white text-center rounded-md p-2 w-24 border border-slate-600 focus:ring-2 focus:ring-blue-500 focus:outline-none"
			placeholder="e.g., 2018"
		/>
	</div>

	<div class="text-slate-400 font-bold text-2xl pt-4">-</div>

	<!-- "After" Year Input -->
	<div class="flex flex-col items-center">
		<label for="end-year" class="text-xs font-bold text-slate-300 mb-1">AFTER</label>
		<input 
			id="end-year"
			type="number"
			bind:value={endYear}
			class="bg-slate-900 text-white text-center rounded-md p-2 w-24 border border-slate-600 focus:ring-2 focus:ring-blue-500 focus:outline-none"
			placeholder="e.g., 2023"
		/>
	</div>
</div>