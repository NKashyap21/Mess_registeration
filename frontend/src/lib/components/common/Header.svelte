<script lang="ts">
	import iithLogoLight from '$lib/assets/iith-logo-light.webp';
	import iithLogoDark from '$lib/assets/iith-logo-dark.webp';
	import darkDarkSvg from '$lib/assets/dark-dark.svg';
	import darkSvg from '$lib/assets/dark.svg';
	import { colorScheme } from '$lib/state/dark.svelte';
	import googleIcon from '$lib/assets/google.webp';
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { goto } from '$app/navigation';
	import CustomDropdown from './CustomDropdown.svelte';
	import { Avatar } from 'bits-ui';
	let { userData }: { userData: any } = $props();
</script>

<header class="flex w-full items-center px-6 pt-4 lg:px-10 lg:pt-8 2xl:px-16 2xl:pt-6">
	<img src={iithLogoDark} class="w-10 lg:w-16 dark:hidden" alt="iith logo" />
	<img src={iithLogoLight} class="hidden w-10 lg:w-16 dark:block" alt="iith logo" />
	<div class="ml-auto flex flex-row items-center gap-x-4 lg:gap-x-8">
		<button
			onclick={() => {
				colorScheme.dark = !colorScheme.dark;
				localStorage.setItem('colorDark', colorScheme.dark ? '1' : '0');
			}}
			class="h-fit hover:cursor-pointer"
		>
			<img src={darkDarkSvg} class="hidden max-lg:w-6 dark:block" alt="dark svg" />
			<img src={darkSvg} class="size max-lg:w-6 dark:hidden" alt="dark svg" />
		</button>

		{#if page.url.pathname != '/login'}
			<CustomDropdown>
				{#snippet trigger()}
					<Avatar.Root
						class="size-10 rounded-full hover:scale-105 hover:cursor-pointer active:brightness-95 lg:size-16"
					>
						<Avatar.Image class="rounded-full" src={userData?.profile_pic} alt="Avatar" />
						<Avatar.Fallback>IH</Avatar.Fallback>
					</Avatar.Root>
				{/snippet}
				{#snippet content()}
					<Avatar.Root class="size-8 rounded-full lg:size-12">
						<Avatar.Image class="rounded-full" src={userData?.profile_pic} alt="Avatar" />
						<Avatar.Fallback>IH</Avatar.Fallback>
					</Avatar.Root>
					<p class="mt-4 text-base font-medium lg:text-lg">{userData.name}</p>
					<p class="font-semibold">{(userData.roll_number as string).toUpperCase()}</p>
					<hr class="my-2 w-3/4 lg:my-4" />
					<button
						onclick={() => {
							fetch(PUBLIC_API_URL + '/logout', {
								method: 'POST',
								credentials: 'include'
							}).finally(() => {
								goto('/login');
							});
						}}
						class="flex flex-row items-center gap-x-4 rounded-md bg-white px-3 py-1.5 text-lg font-semibold text-custom-black shadow-xs shadow-custom-dark-grey hover:brightness-95 active:scale-[99%] active:brightness-90 lg:px-4 lg:py-2 dark:bg-custom-black dark:text-white"
					>
						<img src={googleIcon} class=" size-6 lg:size-8" alt="google icon" />
						<p class="leading-none max-lg:text-sm">Sign Out</p>
					</button>
				{/snippet}
			</CustomDropdown>
		{/if}
	</div>
</header>
