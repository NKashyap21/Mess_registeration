<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Button from '$lib/components/common/Button.svelte';
	import Progress from '$lib/components/common/Progress.svelte';
	let { data } = $props();
</script>

<section class="flex flex-col">
	<div class="text-2xl font-bold">Registration Stats</div>
	<div class="flex items-center gap-x-8">
		<div class="flex flex-col">
			<div class="mt-12 flex flex-col gap-y-8">
				{#each Object.keys(data.messStats.stats) as messName}
					<section class="flex flex-row">
						<h2 class="mr-8 w-[4rem] text-2xl font-bold">{messName}</h2>
						<div class="flex flex-col gap-y-4">
							{#each Object.keys(data.messStats.stats[messName]) as hall, idx}
								<div class="flex items-center">
									<h3 class="w-[4rem] text-xl">{hall}</h3>
									<Progress
										outerClass="min-w-[30rem]"
										innerClass={idx % 2 == 0
											? 'bg-[#919191] dark:bg-[#B4B4B4]'
											: 'bg-custom-light-orange'}
										value={1}
									/>
									<h3 class="px-6 text-xl font-medium">
										{data.messStats.stats[messName][hall].count}/{data.messStats.stats[messName][
											hall
										].capacity}
									</h3>
								</div>
							{/each}
						</div>
					</section>
				{/each}
			</div>
			<!-- <div class="mt-12 grid grid-cols-2 gap-x-4 *:odd:font-medium"> -->
			<!-- 	<p>Next Veg Registration Date:</p> -->
			<!-- 	<p>26 Aug, 5pm</p> -->
			<!-- 	<p>Next Regular Registration Date:</p> -->
			<!-- 	<p>26 Aug, 5pm</p> -->
			<!-- </div> -->
		</div>
		<div class="ml-auto flex flex-col gap-y-6 *:px-8 **:w-full **:text-lg">
			<Button>Import New List</Button>
			<Button
				onclick={() => {
					if (data.registrationState.regular) {
						fetch(PUBLIC_API_URL + '/office/end-registration/reg', {
							method: 'POST',
							credentials: 'include'
						}).then(async (res) => {
							alert(await res.json());
						});
					} else {
						fetch(PUBLIC_API_URL + '/office/start-registration/reg', {
							method: 'POST',
							credentials: 'include'
						}).then(async (res) => {
							alert(await res.json());
						});
					}
					invalidateAll();
				}}>{data.registrationState.regular ? 'Stop' : 'Start'} Regular Registration</Button
			>
			<Button
				onclick={() => {
					if (data.registrationState.veg) {
						fetch(PUBLIC_API_URL + '/office/end-registration/veg', {
							method: 'POST',
							credentials: 'include'
						}).then(async (res) => {
							alert(await res.json());
						});
					} else {
						fetch(PUBLIC_API_URL + '/office/start-registration/veg', {
							method: 'POST',
							credentials: 'include'
						}).then(async (res) => {
							alert(await res.json());
						});
					}
					invalidateAll();
				}}>{data.registrationState.veg ? 'Stop' : 'Start'} Veg Registration</Button
			>
			<!-- <Button>Enable Latest Registration</Button> -->
		</div>
	</div>
</section>
