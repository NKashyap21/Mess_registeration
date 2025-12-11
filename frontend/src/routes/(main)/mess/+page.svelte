<script lang="ts">
	import Button from '$lib/components/common/Button.svelte';
	import CustomSelect from '$lib/components/common/CustomSelect.svelte';
	import MainPanelHeader from '$lib/components/common/MainPanelHeader.svelte';
	import successSvg from '$lib/assets/success.svg';
	import failedSvg from '$lib/assets/failed.svg';
	import Modal from '$lib/components/common/Modal.svelte';
	import { Dialog } from 'bits-ui';
	import { PUBLIC_API_URL } from '$env/static/public';
	let messTimeValue = $state(
		(() => {
			let currentHour = new Date(new Date()).getHours();
			if (currentHour <= 10) {
				return 'Breakfast';
			} else if (currentHour <= 15) {
				return 'Lunch';
			} else if (currentHour <= 19) {
				return 'Snack';
			} else {
				return 'Dinner';
			}
		})()
	);
	// const messSelectItems = [
	// 	{ label: 'Mess A', value: 'Mess A' },
	// 	{ label: 'Mess B', value: 'Mess B' }
	// ];
	const messTimeItems = [
		{ label: 'Breakfast', value: 'Breakfast' },
		{ label: 'Lunch', value: 'Lunch' },
		{ label: 'Snack', value: 'Snack' },
		{ label: 'Dinner', value: 'Dinner' }
	];
	let success = $state(true);
	let rollNo = $state('');
	let studentName = $state('');
	let registeredMess = $state('');
	let extraText = $state('');
</script>

<MainPanelHeader>Dining Scan</MainPanelHeader>

<div class="mt-16 flex flex-row px-32">
	<section class="flex flex-col">
		<div class="mb-36 grid grid-cols-2 items-center gap-y-6 text-xl font-medium">
			<p>Roll No.:</p>
			<input
				bind:value={rollNo}
				onkeydown={async (ev) => {
					if (ev.key == 'Enter') {
						const res = await fetch(
							PUBLIC_API_URL +
								'/messStaff/scanning?roll_no=' +
								rollNo.trim() +
								'&meal=' +
								messTimeValue,
							{
								method: 'GET',
								credentials: 'include'
							}
						);

						const jsondata = await res.json();
						console.log(jsondata);
						if (res.status == 200) {
							success = true;
						} else {
							success = false;
						}
						if (Object.keys(jsondata).includes('data')) {
							if (jsondata['data'] != null) {
								rollNo = jsondata['data']['user']['roll_no'];
								studentName = jsondata['data']['user']['name'];
								const messId = jsondata['data']['user']['mess'];
								if (messId == 1) {
									registeredMess = 'Mess A - LDH';
								} else if (messId == 2) {
									registeredMess = 'Mess A - UDH';
								} else if (messId == 3) {
									registeredMess = 'Mess B - LDH';
								} else if (messId == 4) {
									registeredMess = 'Mess B - UDH';
								}
							}
						}
						extraText = jsondata.message;
						rollNo = '';
					}
				}}
				type="text"
				class="inline-block w-[12rem] rounded-md px-4 py-3 leading-0 outline outline-custom-dark-grey dark:bg-custom-mid-grey dark:outline-custom-light-grey"
			/>
			<p>Student Name:</p>
			<p>{studentName}</p>
			<p>Registered Mess</p>
			<p>{registeredMess}</p>
		</div>
		<div class="flex gap-x-8">
			<!-- <CustomSelect value={messValue} items={messSelectItems} /> -->
			<CustomSelect
				name="meal"
				bind:value={messTimeValue}
				widthClass="!w-[8rem]"
				items={messTimeItems}
			/>
		</div>
	</section>
	<section class="flex flex-col items-center pl-56">
		<div class="flex *:size-16">
			{#if success}
				<img src={successSvg} alt="success svg" />
			{:else}
				<img src={failedSvg} alt="failed svg" />
			{/if}
		</div>
		<div class="flex h-fit flex-col items-center">
			<p class="mt-4 text-xl font-bold">{success ? 'Successful' : 'Failed'}</p>
			<p class="mt-2 text-lg">{extraText}</p>
		</div>

		<Modal buttonText="Flag ID Card" class="mt-auto !text-xl">
			<div class="flex flex-col items-center px-8 py-4">
				<p class="mb-10 text-3xl font-medium">Flag ID?</p>

				<CustomSelect value={'Reason'} items={messTimeItems} />
				<Dialog.Close>
					<Button class="mt-8 !px-16 !py-1.5">Yes</Button>
				</Dialog.Close>
			</div>
		</Modal>
	</section>
</div>
