<script lang="ts">
	import type { HTMLButtonAttributes } from 'svelte/elements';
	let {
		text,
		active,
		class: className,
		...rest
	}: { active: boolean; text: string } & HTMLButtonAttributes = $props();
</script>

<button
	class="rounded-t-2xl text-xl 2xl:text-2xl {className} px-8 py-4 text-center leading-0 font-bold text-nowrap select-none hover:cursor-pointer hover:brightness-95 2xl:px-12 2xl:py-6 {active
		? 'bg-white dark:bg-custom-dark-grey'
		: 'bg-custom-light-grey dark:bg-custom-light-grey'}"
	{...rest}
>
	{text}
</button>
