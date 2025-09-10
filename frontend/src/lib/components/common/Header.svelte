<script lang="ts">
	import iithLogoLight from '$lib/assets/iith-logo-light.webp';
	import iithLogoDark from '$lib/assets/iith-logo-dark.webp';
	import darkSvg from '$lib/assets/dark.svg';
	import { colorScheme } from '$lib/state/dark.svelte';
</script>

<header class="flex w-full items-center px-12 pt-6">
	<img src={iithLogoDark} class="w-16 dark:hidden" alt="iith logo" />
	<img src={iithLogoLight} class="hidden w-16 dark:block" alt="iith logo" />
	<button
		onclick={() => {
			colorScheme.dark = !colorScheme.dark;
		}}
		class="ml-auto hover:cursor-pointer"
	>
		<img src={darkSvg} alt="dark svg" />
	</button>
</header>
