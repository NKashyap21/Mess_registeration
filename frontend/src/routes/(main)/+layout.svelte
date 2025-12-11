<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import NavTabs from '$lib/components/admin/NavTabs.svelte';

	let { children } = $props();
</script>

<div
	class="rounded-xl bg-[#f8f8f8] p-8 shadow-md 2xl:p-12 {page.url.pathname.includes('admin')
		? 'min-w-3/4 rounded-tl-none'
		: ''} dark:bg-custom-dark-grey dark:text-white"
>
	{#if page.url.pathname.includes('admin')}
		<div class="relative">
			<div class="absolute -top-8 -left-8 flex -translate-y-full gap-x-4 2xl:-top-12 2xl:-left-12">
				<NavTabs
					onclick={() => {
						goto(resolve('/admin/registration'));
					}}
					active={page.url.pathname == resolve('/admin/registration')}
					text="Registration"
				/>
				<NavTabs
					onclick={() => {
						goto(resolve('/admin/add'));
					}}
					active={page.url.pathname == resolve('/admin/add')}
					text="Add Registration"
				/>
				<NavTabs
					onclick={() => {
						goto(resolve('/admin/swapping'));
					}}
					active={page.url.pathname == resolve('/admin/swapping')}
					text="Swapping"
				/>
				<NavTabs
					onclick={() => {
						goto(resolve('/admin/flags'));
					}}
					active={page.url.pathname == resolve('/admin/flags')}
					text="Flags"
				/>
				<NavTabs
					onclick={() => {
						goto(resolve('/admin/student'));
					}}
					active={page.url.pathname == resolve('/admin/student')}
					text="Student Search"
				/>
			</div>
		</div>
	{/if}
	{@render children?.()}
</div>
