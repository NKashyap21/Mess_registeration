<script>
	import { PUBLIC_API_URL } from '$env/static/public';
	import Button from '$lib/components/common/Button.svelte';
	import CustomScrollbar from '$lib/components/common/CustomScrollbar.svelte';
	import Modal from '$lib/components/common/Modal.svelte';

	let addToListConfirmMondal = $state(false);
	let { data } = $props();
</script>

<div class="flex px-16 py-6 dark:text-white">
	<h1 class="absolute top-12 right-0 left-0 mx-auto w-fit text-[2.5rem] font-semibold">
		Mess Swapping
	</h1>

	<div class="flex gap-x-10">
		<div class="flex flex-col items-center gap-y-4">
			<h1 class="text-4xl font-medium">Mess A to Mess B</h1>
			<CustomScrollbar>
				<div
					class="relative grid h-[16rem] w-[26rem] auto-cols-max auto-rows-min grid-cols-3 gap-y-1 text-center text-nowrap *:font-medium"
				>
					<div>Sl no.</div>
					<div class="mb-1">Name</div>
					<div>Roll No.</div>
					{#each data.swapData.filter((val) => val.direction == 'A to B') as swap, idx (swap)}
						<div>{idx + 1}</div>
						<div class={swap.name == data.user.name ? 'text-custom-orange' : ''}>
							{swap.name}{swap.name == data.user.name ? ' (You)' : ''}
						</div>
						<div class={swap.name == data.user.name ? 'text-custom-orange' : ''}>
							{swap.roll_no}
						</div>
					{/each}
				</div>
			</CustomScrollbar>
		</div>
		<div class="flex flex-col items-center gap-y-4">
			<h1 class="text-4xl font-medium">Mess B to Mess A</h1>
			<CustomScrollbar>
				<div
					class="relative grid h-[16rem] w-[26rem] auto-rows-min grid-cols-3 gap-y-1 text-center *:font-medium"
				>
					<div>Sl no.</div>
					<div class="mb-1">Name</div>
					<div>Roll No.</div>

					{#each data.swapData.filter((val) => val.direction == 'B to A') as swap, idx (swap)}
						<div>{idx + 1}</div>
						<div class={swap.name == data.user.name ? 'text-custom-orange' : ''}>
							{swap.name}{swap.name == data.user.name ? ' (You)' : ''}
						</div>
						<div class={swap.name == data.user.name ? 'text-custom-orange' : ''}>
							{swap.roll_no}
						</div>
					{/each}
				</div>
			</CustomScrollbar>
		</div>
	</div>
</div>

<div class="absolute right-0 bottom-40 left-0 mx-auto flex items-center justify-center">
	<Modal buttonText={'Add to List'} bind:open={addToListConfirmMondal}>
		<div class="flex flex-col items-center text-nowrap">
			<h2 class="text-3xl">Swap Mess</h2>
			<p class="mt-4 text-lg">Are you sure you want to swap your current mess to the other?</p>
			<div class="mt-8 flex gap-x-4 self-end">
				<Button
					onclick={() => {
						addToListConfirmMondal = false;
					}}>No</Button
				>
				<Button
					onclick={() => {
						fetch(PUBLIC_API_URL + '/students/createSwap', {
							method: 'POST',
							credentials: 'include'
						})
							.then((res) => {
								console.log(res);
								if (res.status == 200 || res.status == 201) {
									res.json().then((data) => {
										console.log(data);
										alert(data.message + '\n' + 'Swap created from ' + data.data.direction);
									});
								} else {
									res.json().then((data) => {
										alert(data.error);
									});
								}
							})
							.catch((e) => {
								console.error(e);
								alert('Failed to add swap request');
							})
							.finally(() => {
								addToListConfirmMondal = false;
							});
					}}>Yes</Button
				>
			</div>
		</div>
	</Modal>
</div>
