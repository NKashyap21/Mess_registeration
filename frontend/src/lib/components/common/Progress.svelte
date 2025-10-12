<script lang="ts">
	import type { ClassValue } from 'svelte/elements';

	let {
		outerClass,
		innerClass,
		value
	}: { outerClass?: ClassValue; innerClass?: ClassValue; value: number } = $props();
</script>

<div class="h-5 overflow-hidden rounded-full {outerClass} bg-[#E0E0E0] p-1 dark:bg-black">
	<div class="h-full rounded-full {innerClass}" style="width: {value}%;"></div>
</div>
