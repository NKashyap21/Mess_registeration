<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/common/Header.svelte';
	import { colorScheme } from '$lib/state/dark.svelte';

	let { children, data } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div
	class="flex h-full transition-colors {colorScheme.dark
		? 'dark'
		: ''} w-full flex-col bg-custom-lighter-grey dark:bg-custom-black"
>
	<Header userData={data.user} />
	<main class="flex h-full w-full grow flex-col items-center justify-center **:transition-colors">
		{@render children?.()}
	</main>
	<div class="mb-10 self-center text-custom-light-grey">Brought to you with 🧡💜🧡 by Lambda</div>
</div>
