<script lang="ts">
	import diningSvg from '$lib/assets/dining.svg';
	import diningDarkSvg from '$lib/assets/dining-dark.svg';
	let { children } = $props();
</script>

<div class="flex w-full flex-col items-center justify-center pb-3 2xl:pb-4">
	<img src={diningSvg} alt="dining icon" class="max-2xl:size-12 dark:hidden" />
	<img src={diningDarkSvg} alt="dining icon" class="hidden max-2xl:size-12 dark:block" />
	<h3 class="mt-3 text-xl font-medium text-custom-black 2xl:text-2xl dark:text-white">
		{@render children()}
	</h3>
</div>
