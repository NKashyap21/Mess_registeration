<script lang="ts">
	import type { PageProps } from './$types';
	import GoogleSignInBtn from '$lib/components/GoogleSignInBtn.svelte';

	let { data }: PageProps = $props();
	console.log(data.registeredMess);
</script>

<main class="p-8">
	<section class="flex items-center">
		<div class="flex flex-col">
			<h1 class="text-4xl font-bold">Mess Registration Portal</h1>
			<h4 class=" opacity-50">Indian Institute of Technology - Hyderabad</h4>
		</div>
		<div class="ml-auto flex items-center gap-x-4">
			{#if data.registrationLive}
				<div class="badge badge-soft p-4 badge-neutral">
					<span class="status animate-pulse status-success"></span>
					Registrations Live
				</div>
			{/if}
			{#if data.session == null}
				<GoogleSignInBtn />
			{:else}
				<div class="flex items-center gap-x-2">
					<img
						class="size-10 rounded-full"
						src={data.session?.user?.image ?? 'https://avatar.iran.liara.run/public'}
						alt="profile pic"
					/>
					<p>{data.session?.user?.name ?? 'User'}</p>
				</div>
			{/if}
		</div>
	</section>

	{#each data.messData as mess}
		<section class="mt-12 flex flex-col gap-y-3">
			<h3 class="text-xl">
				{mess.name}
				{#if data.registeredMess[0].name == mess.name}
					<span class="text-lg opacity-75">( Registered for this mess )</span>
				{/if}
			</h3>
			<div class="stats shadow">
				<div class="stat">
					<p class="stat-value">{mess.total_registrants}</p>
					<p class="stat-title">Total Registrants</p>
				</div>
				<div class="stat">
					<p class="stat-value">{mess.total_capacity - mess.total_registrants}</p>
					<p class="stat-title">Available Slots</p>
				</div>
				<div class="stat">
					<p class="stat-value">{mess.total_capacity}</p>
					<p class="stat-title">Total Slots</p>
				</div>
			</div>
			<button
				disabled={data.session == null || !data.registrationLive}
				class="btn w-fit self-end {data.session != null ? 'btn-primary' : 'btn-warning'}"
				hidden={data.registeredMess.length != 0}
			>
				{#if data.registrationLive}
					{#if data.session != null}
						Register for this mess
					{:else}
						Log in to register
					{/if}
				{:else}
					Registration has not started yet
				{/if}
			</button>
		</section>
	{/each}
</main>
