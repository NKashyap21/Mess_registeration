<script lang="ts">
	import type { SearchEmailResponse } from '$lib/types/interface';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let registrationLive = $state(data.registrationLive);
	let formEmail = $state('');

	let searchUserRegistered = $state(false);
	let searchUserMess = $state(-1);
</script>

<main class="flex flex-col p-8">
	<section class="mb-10 flex items-center">
		<div class="flex flex-col">
			<h1 class="text-4xl font-bold">Mess Admin Portal</h1>
			<h4 class=" opacity-50">Indian Institute of Technology - Hyderabad</h4>
		</div>
		<div class="ml-auto flex items-center gap-x-4">
			{#if data.session != null}
				<div class="flex items-center gap-x-2">
					{@debug data}
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

	<div class="flex">
		<section class="flex gap-x-4">
			<div class="flex flex-col">
				{#each data.messData as mess}
					<h4 class="mb-2">{mess.name}</h4>
					<div class="stats mb-5 border border-base-300">
						<div class="stat">
							<div class="stat-title">Total Capacity</div>
							<div class="stat-value">{mess.total_capacity}</div>
						</div>
						<div class="stat">
							<div class="stat-title">Total Registered</div>
							<div class="stat-value">{mess.total_registrants}</div>
						</div>
					</div>
				{/each}
			</div>
			<div class="flex flex-col gap-y-6">
				<div class="stats h-min border border-base-300">
					<div class="stat">
						<div class="stat-title">Turn on Registrations</div>
						<div class="stat-value">
							{#if registrationLive}
								Live
							{:else}
								Offline
							{/if}
						</div>
						<div class="stat-actions grow">
							{#if registrationLive}
								<button
									class="btn btn-sm btn-warning"
									onclick={() => {
										registrationLive = false;
									}}>Turn Off</button
								>
							{:else}
								<button
									class="btn btn-sm btn-primary"
									onclick={() => {
										registrationLive = true;
									}}>Turn On</button
								>
							{/if}
						</div>
					</div>
				</div>
				<div class="card border border-base-300">
					<div class="card-body">
						<p class="card-title">Swap Mess</p>
						<div class="card-actions justify-end">
							<button class="btn btn-sm btn-primary">Swap</button>
						</div>
					</div>
				</div>
			</div>
		</section>

		<section class="ml-4 flex w-full flex-col items-center rounded-box border border-base-300 p-8">
			<div class="flex flex-col items-center gap-y-2">
				<h4 class="text-lg">Enter email id to search</h4>
				<div class="join">
					<div>
						<label class="validator input join-item">
							<svg
								class="h-[1em] opacity-50"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 24 24"
							>
								<g
									stroke-linejoin="round"
									stroke-linecap="round"
									stroke-width="2.5"
									fill="none"
									stroke="currentColor"
								>
									<rect width="20" height="16" x="2" y="4" rx="2"></rect>
									<path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"></path>
								</g>
							</svg>
							<input
								bind:value={formEmail}
								type="email"
								name="email"
								placeholder="mail@iith.ac.in"
								class="w-[15rem]"
								required
							/>
						</label>
						<div class="validator-hint hidden">Enter valid email address</div>
					</div>
					<button
						onclick={async () => {
							const params = new URLSearchParams({ email: formEmail });
							const res = await fetch('admin/search?' + params.toString());
							const data = (await res.json()) as SearchEmailResponse;
							searchUserMess = data.mess_id;
							searchUserRegistered = data.mess_id != -1;
						}}
						class="btn join-item btn-primary">Search</button
					>
				</div>
			</div>
			<div class="stats mt-8 border border-base-300">
				<div class="stat">
					<p class="stat-title">Registered</p>
					<p class="stat-value">
						{#if searchUserRegistered}
							Yes
						{:else}
							No
						{/if}
					</p>
				</div>

				<div class="stat">
					<p class="stat-title">Mess</p>
					<p class="stat-value">
						{#if searchUserMess != -1}
							{data.messData.find((mess) => mess.id == searchUserMess)?.name}
						{:else}
							No Mess
						{/if}
					</p>
				</div>
			</div>
		</section>
	</div>
</main>
