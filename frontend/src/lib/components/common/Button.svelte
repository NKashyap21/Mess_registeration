<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { ClassValue, HTMLButtonAttributes } from 'svelte/elements';

	let {
		children,
		class: className,
		...rest
	}: HTMLButtonAttributes & { children: Snippet; class?: ClassValue } = $props();
</script>

<button
	{...rest}
	class="w-fit rounded-md bg-custom-off-white px-10 py-2.5 text-xl font-medium shadow-xs shadow-custom-dark-grey outline outline-custom-red hover:brightness-95 active:scale-[99%] active:brightness-90 dark:bg-custom-mid-grey dark:text-white dark:outline-custom-light-grey {className}"
>
	{@render children()}
</button>
