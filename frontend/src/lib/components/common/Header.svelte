<script lang="ts">
	import iithLogoLight from '$lib/assets/iith-logo-light.webp';
	import iithLogoDark from '$lib/assets/iith-logo-dark.webp';
	import darkDarkSvg from '$lib/assets/dark-dark.svg';
	import darkSvg from '$lib/assets/dark.svg';
	import { colorScheme } from '$lib/state/dark.svelte';
	import googleIcon from '$lib/assets/google.webp';
	import { page } from '$app/state';
</script>

<header class="flex w-full items-center px-12 pt-6">
	<img src={iithLogoDark} class="w-16 dark:hidden" alt="iith logo" />
	<img src={iithLogoLight} class="hidden w-16 dark:block" alt="iith logo" />
	<div class="ml-auto flex flex-row items-center gap-x-8">
		<button
			onclick={() => {
				colorScheme.dark = !colorScheme.dark;
				localStorage.setItem('colorDark', colorScheme.dark ? '1' : '0');
			}}
			class="h-fit hover:cursor-pointer"
		>
			<img src={darkDarkSvg} class="hidden dark:block" alt="dark svg" />
			<img src={darkSvg} class="size dark:hidden" alt="dark svg" />
		</button>

		{#if page.url.pathname != '/login'}
			<button
				class="flex flex-row items-center gap-x-4 rounded-md bg-white px-8 py-5 text-xl font-semibold text-custom-black shadow-xs shadow-custom-dark-grey hover:brightness-95 active:scale-[99%] active:brightness-90 dark:bg-custom-dark-grey dark:text-white"
			>
				<img src={googleIcon} class="size-8" alt="google icon" />
				<p class="leading-none">Sign Out</p>
			</button>
		{/if}
	</div>
</header>
