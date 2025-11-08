<script lang="ts">
	import { goto } from '$app/navigation';
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
						goto('/admin/registration');
					}}
					active={page.url.pathname == '/admin/registration'}
					text="Registration"
				/>
				<NavTabs
					onclick={() => {
						goto('/admin/add');
					}}
					active={page.url.pathname == '/admin/add'}
					text="Add Registration"
				/>
				<NavTabs
					onclick={() => {
						goto('/admin/swapping');
					}}
					active={page.url.pathname == '/admin/swapping'}
					text="Swapping"
				/>
				<NavTabs
					onclick={() => {
						goto('/admin/flags');
					}}
					active={page.url.pathname == '/admin/flags'}
					text="Flags"
				/>
				<NavTabs
					onclick={() => {
						goto('/admin/student');
					}}
					active={page.url.pathname == '/admin/student'}
					text="Student Search"
				/>
			</div>
		</div>
	{/if}
	{@render children?.()}
</div>
