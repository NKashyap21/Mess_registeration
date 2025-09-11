<script lang="ts">
	import diningSvg from '$lib/assets/dining.svg';
	import diningDarkSvg from '$lib/assets/dining-dark.svg';
	let { children } = $props();
</script>

<div class="flex w-full flex-col items-center justify-center pb-4">
	<img src={diningSvg} alt="dining icon" class="dark:hidden" />
	<img src={diningDarkSvg} alt="dining icon" class="hidden dark:block" />
	<h3 class="mt-3 text-2xl font-medium text-custom-black dark:text-white">{@render children()}</h3>
</div>
