<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Button from '$lib/components/common/Button.svelte';
	import CustomScrollbar from '$lib/components/common/CustomScrollbar.svelte';
	import Modal from '$lib/components/common/Modal.svelte';

	let createSwapModal = $state(false);
	let deleteSwapModal = $state(false);
	let acceptSwapModal = $state(false);
	let { data } = $props();

	let swapType = $state<'public' | 'friend'>('public');
	let swapPassword = $state('');
	let isCreating = $state(false);

	let acceptPassword = $state('');
	let selectedSwap: { user_id: number; type: string; name: string; direction: string } | null =
		$state(null);
	let isAccepting = $state(false);

	let isDeleting = $state(false);

	async function handleCreateSwap() {
		if (swapType === 'friend' && (swapPassword.length < 6 || swapPassword.length > 100)) {
			alert('Password must be between 6 and 100 characters');
			return;
		}

		isCreating = true;
		try {
			const response = await fetch(PUBLIC_API_URL + '/students/createSwap', {
				method: 'POST',
				credentials: 'include',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					type: swapType,
					password: swapType === 'friend' ? swapPassword : 'public_swap'
				})
			});

			const result = await response.json();
			if (response.ok) {
				alert(result.message + '\nSwap created: ' + result.data.direction);
				createSwapModal = false;
				swapType = 'public';
				swapPassword = '';
				invalidateAll();
			} else {
				alert(result.error ?? 'Failed to create swap request');
			}
		} catch (err) {
			console.error(err);
			alert('Failed to create swap request');
		} finally {
			isCreating = false;
		}
	}

	async function handleDeleteSwap() {
		isDeleting = true;
		try {
			const response = await fetch(PUBLIC_API_URL + '/students/deleteSwap', {
				method: 'DELETE',
				credentials: 'include'
			});

			const result = await response.json();
			if (response.ok) {
				alert(result.message ?? 'Swap request deleted');
				deleteSwapModal = false;
				invalidateAll();
			} else {
				alert(result.error ?? 'Failed to delete swap request');
			}
		} catch (err) {
			console.error(err);
			alert('Failed to delete swap request');
		} finally {
			isDeleting = false;
		}
	}

	async function handleAcceptSwap() {
		if (!selectedSwap) return;

		if (selectedSwap.type === 'friend' && acceptPassword.length < 6) {
			alert('Please enter the password');
			return;
		}

		isAccepting = true;
		try {
			console.log(selectedSwap);

			const response = await fetch(PUBLIC_API_URL + '/students/acceptSwap', {
				method: 'POST',
				credentials: 'include',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					type: selectedSwap.type,
					user_id: selectedSwap.user_id,
					password: selectedSwap.type === 'friend' ? acceptPassword : ''
				})
			});

			const result = await response.json();
			if (response.ok) {
				alert(result.message ?? 'Swap accepted successfully');
				acceptSwapModal = false;
				selectedSwap = null;
				acceptPassword = '';
				invalidateAll();
			} else {
				alert(result.error ?? 'Failed to accept swap request');
			}
		} catch (err) {
			console.error(err);
			alert('Failed to accept swap request');
		} finally {
			isAccepting = false;
		}
	}

	function openAcceptModal(swap: any) {
		selectedSwap = {
			user_id: swap.user_id,
			type: swap.type,
			name: swap.name,
			direction: swap.direction
		};
		acceptPassword = '';
		acceptSwapModal = true;
	}

	console.log(data);
</script>

<div class="flex flex-col px-4 py-6 sm:px-8 md:px-16 dark:text-white">
	<h1 class="mb-6 text-center text-2xl font-semibold sm:text-3xl md:text-[2.5rem]">
		Mess Swapping
	</h1>

	{#if data.mySwap}
		<div
			class="mx-auto mb-6 w-full max-w-md rounded-lg border border-gray-300 bg-gray-50 p-4 sm:mb-8 sm:p-6 dark:border-gray-600 dark:bg-gray-900"
		>
			<div class="flex flex-col items-center gap-y-3 sm:gap-y-4">
				<h2 class="text-base font-semibold sm:text-lg">Your Active Swap Request</h2>
				<div class="flex gap-x-6 text-center sm:gap-x-8">
					<div class="flex flex-col">
						<span class="text-xs text-gray-500 dark:text-gray-400">Type</span>
						<span class="text-sm font-medium capitalize"
							>{data.mySwap.type === 'friend' ? 'Private' : 'Public'}</span
						>
					</div>
					<div class="flex flex-col">
						<span class="text-xs text-gray-500 dark:text-gray-400">Status</span>
						<span class="text-sm font-medium">Active</span>
					</div>
				</div>
				<Button onclick={() => (deleteSwapModal = true)}>Cancel Request</Button>
			</div>
		</div>
	{/if}

	<div class="flex w-full justify-center">
		<div class="flex w-full max-w-2xl flex-col items-center gap-y-3 sm:gap-y-4">
			<h2 class="text-center text-2xl font-medium sm:text-3xl md:text-4xl">
				Private Swap Requests
			</h2>
			<p class="text-xs text-gray-500 sm:text-sm dark:text-gray-400">
				Click on a request to accept it
			</p>
			<CustomScrollbar>
				<div class="flex h-[16rem] w-full flex-col gap-y-2 pr-2 sm:h-[20rem] sm:pr-4">
					{#each data.swapData.filter((s: any) => s.type === 'friend' && s.name !== data.user.name) as swap, idx (swap)}
						<button
							class="flex items-center justify-between rounded-lg border border-gray-300 bg-gray-50 px-3 py-2 text-left transition-colors hover:bg-gray-100 sm:px-4 sm:py-3 dark:border-gray-600 dark:bg-gray-900 dark:hover:bg-gray-800"
							onclick={() => openAcceptModal(swap)}
						>
							<div class="flex min-w-0 flex-1 items-center gap-x-2 sm:gap-x-4">
								<span class="shrink-0 text-xs text-gray-500 sm:text-sm">{idx + 1}.</span>
								<div class="flex min-w-0 flex-col">
									<span class="truncate text-sm font-medium sm:text-base" title={swap.name}
										>{swap.name}</span
									>
									<span class="text-xs text-gray-500">{swap.direction}</span>
								</div>
							</div>
							<span class="ml-2 shrink-0 text-xs text-custom-orange sm:text-sm">Accept →</span>
						</button>
					{:else}
						<div class="flex h-full items-center justify-center text-gray-500 text-sm">
							No private swap requests available
						</div>
					{/each}
				</div>
			</CustomScrollbar>
		</div>
	</div>

	{#if !data.mySwap}
		<div class="mt-6 flex items-center justify-center sm:mt-8">
			<Modal buttonText="Create Swap Request" bind:open={createSwapModal}>
				<div
					class="flex w-full max-w-[90vw] min-w-0 flex-col gap-y-4 px-4 py-4 sm:min-w-[400px] sm:gap-y-6 sm:px-8 sm:py-6"
				>
					<h2 class="text-xl font-semibold sm:text-2xl">Create Swap Request</h2>

					<div class="flex flex-col gap-y-4">
						<div class="flex flex-col gap-y-2">
							<label class="text-base font-medium sm:text-lg">Swap Type</label>
							<div class="flex flex-wrap gap-x-4 gap-y-2">
								<label class="flex cursor-pointer items-center gap-x-2">
									<input
										type="radio"
										name="swapType"
										value="public"
										bind:group={swapType}
										class="h-4 w-4"
									/>
									<span class="text-sm sm:text-base">Public</span>
								</label>
								<label class="flex cursor-pointer items-center gap-x-2">
									<input
										type="radio"
										name="swapType"
										value="friend"
										bind:group={swapType}
										class="h-4 w-4"
									/>
									<span class="text-sm sm:text-base">Private (Friend)</span>
								</label>
							</div>
							<p class="text-xs text-gray-500 sm:text-sm">
								{swapType === 'public'
									? 'Public swaps are visible to everyone and can be accepted by anyone.'
									: 'Private swaps require a password. Share it with your friend to swap.'}
							</p>
						</div>

						{#if swapType === 'friend'}
							<div class="flex flex-col gap-y-2">
								<label for="swapPassword" class="text-base font-medium sm:text-lg">
									Password (share with friend)
								</label>
								<input
									id="swapPassword"
									type="password"
									bind:value={swapPassword}
									placeholder="Enter password (6-100 characters)"
									class="rounded-lg border border-gray-300 px-3 py-2 text-sm sm:px-4 sm:text-base dark:border-gray-600 dark:bg-gray-800 dark:text-white"
								/>
							</div>
						{/if}
					</div>

					<div class="flex flex-col-reverse gap-2 sm:ml-auto sm:flex-row sm:gap-x-4">
						<Button
							onclick={() => {
								createSwapModal = false;
								swapType = 'public';
								swapPassword = '';
							}}>Cancel</Button
						>
						<Button onclick={handleCreateSwap} disabled={isCreating}>
							{isCreating ? 'Creating...' : 'Create Request'}
						</Button>
					</div>
				</div>
			</Modal>
		</div>
	{/if}
</div>

{#if acceptSwapModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
		<div class="w-full max-w-md rounded-lg bg-white p-4 sm:p-6 dark:bg-gray-800 dark:text-white">
			<h2 class="mb-4 text-xl font-semibold sm:text-2xl">Accept Private Swap</h2>

			{#if selectedSwap}
				<div class="mb-6 flex flex-col gap-y-4">
					<div
						class="rounded-lg border border-gray-200 bg-gray-50 p-3 sm:p-4 dark:border-gray-700 dark:bg-gray-900"
					>
						<div class="grid grid-cols-2 gap-x-3 gap-y-2 text-sm sm:gap-x-4">
							<span class="text-gray-500">From:</span>
							<span class="truncate font-medium">{selectedSwap.name}</span>
							<span class="text-gray-500">Direction:</span>
							<span class="font-medium">{selectedSwap.direction}</span>
						</div>
					</div>

					<div class="flex flex-col gap-y-2">
						<label for="acceptPassword" class="text-sm font-medium">Enter Password</label>
						<input
							id="acceptPassword"
							type="password"
							bind:value={acceptPassword}
							placeholder="Enter the password shared by your friend"
							class="rounded-lg border border-gray-300 px-3 py-2 text-sm sm:px-4 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
						/>
						<p class="text-xs text-gray-500">
							Ask your friend for the password they set when creating this swap request.
						</p>
					</div>
				</div>
			{/if}

			<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end sm:gap-x-4">
				<Button
					onclick={() => {
						acceptSwapModal = false;
						selectedSwap = null;
						acceptPassword = '';
					}}>Cancel</Button
				>
				<Button onclick={handleAcceptSwap} disabled={isAccepting}>
					{isAccepting ? 'Swapping...' : 'Swap'}
				</Button>
			</div>
		</div>
	</div>
{/if}

{#if deleteSwapModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
		<div class="w-full max-w-md rounded-lg bg-white p-4 sm:p-6 dark:bg-gray-800 dark:text-white">
			<h2 class="mb-4 text-lg font-semibold sm:text-xl">Cancel Swap Request</h2>
			<p class="mb-6 text-sm text-gray-600 sm:text-base dark:text-gray-300">
				Are you sure you want to cancel your swap request? This action cannot be undone.
			</p>
			<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end sm:gap-x-4">
				<Button onclick={() => (deleteSwapModal = false)}>No, Keep It</Button>
				<Button onclick={handleDeleteSwap} disabled={isDeleting}>
					{isDeleting ? 'Cancelling...' : 'Yes, Cancel Request'}
				</Button>
			</div>
		</div>
	</div>
{/if}
