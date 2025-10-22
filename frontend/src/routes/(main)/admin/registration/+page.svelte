<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Button from '$lib/components/common/Button.svelte';
	import Modal from '$lib/components/common/Modal.svelte';
	import Progress from '$lib/components/common/Progress.svelte';
	let { data } = $props();
	let vegConfirm = $state(false);
	let regConfirm = $state(false);
	let applyConfirm = $state(false);
</script>

<section class="flex flex-col">
	<div class="text-2xl font-bold">Registration Stats</div>
	<div class="flex items-center gap-x-8 text-nowrap">
		<div class="mt-8 flex flex-col gap-y-8 *:gap-y-4">
			<div class="flex flex-col">
				<div class="text-xl font-semibold">Current Stats</div>
				<div class="flex flex-col">
					<div class="flex flex-col gap-y-6">
						{#each Object.keys(data.status['current_mess']) as messName}
							<section class="flex flex-row">
								<h2 class=" min-w-[6rem] text-xl font-bold">{messName}</h2>
								<div class="flex flex-col gap-y-2">
									{#each Object.keys(data.status['current_mess'][messName]) as hall, idx}
										<div class="flex items-center">
											<h3 class="w-[8rem] text-xl">{hall}</h3>
											<Progress
												outerClass="min-w-[30rem]"
												innerClass={idx % 2 == 0
													? 'bg-[#919191] dark:bg-[#B4B4B4]'
													: 'bg-custom-light-orange'}
												value={(data.status['current_mess'][messName][hall] * 100) /
													data.status['capacity'][messName][hall]}
											/>
											<h3 class="px-6 text-xl font-medium">
												{data.status['current_mess'][messName][hall]}/{data.status['capacity'][
													messName
												][hall]}
											</h3>
										</div>
									{/each}
								</div>
							</section>
						{/each}
					</div>
				</div>
			</div>

			<div class="flex flex-col">
				<div class="text-xl font-semibold">Upcoming Stats</div>
				<div class="flex flex-col">
					<div class="flex flex-col gap-y-6">
						{#each Object.keys(data.status['upcoming_mess']) as messName}
							<section class="flex flex-row">
								<h2 class="w-[6rem] text-xl font-bold">{messName}</h2>
								<div class="flex flex-col gap-y-2">
									{#each Object.keys(data.status['upcoming_mess'][messName]) as hall, idx}
										<div class="flex items-center">
											<h3 class="w-[8rem] text-xl">{hall}</h3>
											<Progress
												outerClass="min-w-[30rem] "
												innerClass={idx % 2 == 0
													? 'bg-[#919191] dark:bg-[#B4B4B4]'
													: 'bg-custom-light-orange'}
												value={(data.status['upcoming_mess'][messName][hall] * 100) /
													data.status['capacity'][messName][hall]}
											/>
											<h3 class="px-6 text-xl font-medium">
												{data.status['upcoming_mess'][messName][hall]}/{data.status['capacity'][
													messName
												][hall]}
											</h3>
										</div>
									{/each}
								</div>
							</section>
						{/each}
					</div>
				</div>
			</div>
			<!-- <div class="mt-12 grid grid-cols-2 gap-x-4 *:odd:font-medium"> -->
			<!-- 	<p>Next Veg Registration Date:</p> -->
			<!-- 	<p>26 Aug, 5pm</p> -->
			<!-- 	<p>Next Regular Registration Date:</p> -->
			<!-- 	<p>26 Aug, 5pm</p> -->
			<!-- </div> -->
		</div>
		<div class="ml-auto flex flex-col gap-y-6 **:w-full **:text-lg">
			<Button>Import New List</Button>
			<Modal
				bind:open={regConfirm}
				buttonText="{data.status['registration_status']['normal']
					? 'Stop'
					: 'Start'} Regular Registration"
				class=""
			>
				<div class="flex flex-col gap-y-8 px-8 py-6">
					<div class="text-xl text-nowrap">
						Are you sure you want to {data.status['registration_status']['normal']
							? 'Stop'
							: 'Start'} regular registration?
					</div>
					<div class="ml-auto flex gap-x-4 self-end">
						<Button
							class=""
							onclick={() => {
								regConfirm = false;
							}}>No</Button
						>
						<Button
							class=""
							onclick={async () => {
								await fetch(PUBLIC_API_URL + '/office/toggle/reg', {
									method: 'POST',
									credentials: 'include'
								}).then(async (res) => {
									const data = await res.json();
									alert(data.message ?? data.error ?? 'Unknown');
								});
								invalidateAll();
								regConfirm = false;
							}}>Yes</Button
						>
					</div>
				</div>
			</Modal>
			<Modal
				bind:open={vegConfirm}
				buttonText="{data.status['registration_status']['veg'] ? 'Stop' : 'Start'} Veg Registration"
			>
				<div class="flex flex-col gap-y-8 px-8 py-6">
					<div class="text-xl text-nowrap">
						Are you sure you want to {data.status['registration_status']['veg'] ? 'Stop' : 'Start'} veg
						registration?
					</div>
					<div class="ml-auto flex gap-x-4 self-end">
						<Button
							class=""
							onclick={async () => {
								vegConfirm = false;
							}}>No</Button
						>
						<Button
							class=""
							onclick={async () => {
								await fetch(PUBLIC_API_URL + '/office/toggle/veg', {
									method: 'POST',
									credentials: 'include'
								}).then(async (res) => {
									const data = await res.json();
									alert(data.message ?? data.error ?? 'Unknown');
								});
								invalidateAll();
								vegConfirm = false;
							}}>Yes</Button
						>
					</div>
				</div>
			</Modal>

			<Modal bind:open={applyConfirm} buttonText="Apply Upcoming Registration">
				<div class="flex flex-col gap-y-8 px-8 py-6">
					<div class="text-xl text-nowrap">
						Are you sure you want to apply upcoming registration to current registration?
					</div>
					<div class="ml-auto flex gap-x-4 self-end">
						<Button
							class=""
							onclick={async () => {
								applyConfirm = false;
							}}>No</Button
						>
						<Button
							class=""
							onclick={async () => {
								await fetch(PUBLIC_API_URL + '/office/apply-new-registration', {
									method: 'POST',
									credentials: 'include'
								}).then(async (res) => {
									const data = await res.json();
									alert(data.message ?? data.error ?? 'Unknown');
								});
								invalidateAll();
								applyConfirm = false;
							}}>Yes</Button
						>
					</div>
				</div>
			</Modal>
			<!-- <Button>Enable Latest Registration</Button> -->
		</div>
	</div>
</section>
