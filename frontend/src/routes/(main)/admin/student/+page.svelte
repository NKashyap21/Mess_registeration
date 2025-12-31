<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import Button from '$lib/components/common/Button.svelte';
	import Modal from '$lib/components/common/Modal.svelte';
	import CustomSelect from '$lib/components/common/CustomSelect.svelte';

	let rollNo = $state('');
	let searchData = $state<{ Name: string; RollNo: string; Mess: number | string; Email: string }>({
		Name: '',
		Mess: '',
		RollNo: '',
		Email: ''
	});

	let swapMessModalOpen = $state(false);
	let deregisterModalOpen = $state(false);

	let selectedMess = $state('');
	let isSubmitting = $state(false);

	const messOptions = [
		{ label: 'Mess A - LDH', value: '1' },
		{ label: 'Mess A - UDH', value: '2' },
		{ label: 'Mess B - LDH', value: '3' },
		{ label: 'Mess B - UDH', value: '4' },
		{ label: 'Veg Mess', value: '5' }
	];

	function getMessLabel(messNumber: number | string): string {
		if (typeof messNumber !== 'number') return 'Unregistered';
		const mess = messOptions.find((m) => m.value === messNumber.toString());
		return mess ? mess.label : 'Unknown';
	}

	function searchStudent() {
		if (!rollNo.trim()) {
			alert('Please enter a roll number');
			return;
		}

		fetch(PUBLIC_API_URL + '/office/students/' + rollNo, { credentials: 'include' })
			.then((res) => {
				if (res.status == 200) {
					res.json().then((data) => {
						searchData = data;
						if (typeof data.Mess === 'number') {
							selectedMess = data.Mess.toString();
						} else {
							selectedMess = '';
						}
					});
				} else if (res.status == 404) {
					alert('Student not found');
				} else {
					alert('Error fetching student data');
				}
			})
			.catch((e) => {
				console.error(e);
				alert('Error connecting to server');
			});
	}

	async function handleSwapMess() {
		if (!selectedMess) {
			alert('Please select a mess');
			return;
		}

		if (!searchData.RollNo) {
			alert('No student selected');
			return;
		}

		isSubmitting = true;

		try {
			const response = await fetch(PUBLIC_API_URL + '/office/students/', {
				method: 'PUT',
				credentials: 'include',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					roll_no: searchData.RollNo,
					mess: parseInt(selectedMess),
					can_register: true
				})
			});

			if (response.ok) {
				const result = await response.json();
				alert(result.message || 'Mess updated successfully');
				swapMessModalOpen = false;
				searchStudent();
			} else {
				const error = await response.json();
				alert(error.error || 'Failed to update mess');
			}
		} catch (e) {
			console.error(e);
			alert('Error connecting to server');
		} finally {
			isSubmitting = false;
		}
	}

	async function handleDeregister() {}
</script>

<section class="flex w-full flex-col">
	<div class="flex text-2xl">
		<input
			class="rounded-md border border-custom-mid-grey bg-white px-6 py-4"
			placeholder="Roll Number"
			bind:value={rollNo}
			onkeyup={(ev) => {
				if (ev.key == 'Enter') {
					searchStudent();
				}
			}}
		/>
		<Button class="ml-4" onclick={searchStudent}>Search</Button>
		<div class="ml-auto flex items-center gap-x-2 text-custom-orange">
			View Student History
			<svg
				width="11"
				height="19"
				viewBox="0 0 11 19"
				fill="none"
				xmlns="http://www.w3.org/2000/svg"
			>
				<path
					fill-rule="evenodd"
					clip-rule="evenodd"
					d="M10.58 8.47353C10.8489 8.74114 11 9.10406 11 9.48246C11 9.86087 10.8489 10.2238 10.58 10.4914L2.46603 18.5643C2.33372 18.7006 2.17545 18.8093 2.00045 18.8841C1.82546 18.9589 1.63725 18.9983 1.4468 18.9999C1.25635 19.0016 1.06747 18.9655 0.8912 18.8937C0.714926 18.822 0.554779 18.716 0.420106 18.582C0.285433 18.448 0.178929 18.2887 0.106809 18.1133C0.0346899 17.9379 -0.0016008 17.75 5.41561e-05 17.5605C0.00170911 17.371 0.041277 17.1838 0.116449 17.0097C0.191621 16.8356 0.300891 16.6781 0.437884 16.5465L7.53783 9.48246L0.437884 2.41847C0.176609 2.14932 0.0320366 1.78884 0.0353046 1.41467C0.0385726 1.0405 0.189419 0.682569 0.455356 0.417979C0.721293 0.153388 1.08104 0.00330541 1.45712 5.3947e-05C1.8332 -0.00319751 2.19551 0.140643 2.46603 0.400595L10.58 8.47353Z"
					fill="#F05A25"
				/>
			</svg>
		</div>
	</div>

	{#if searchData.RollNo}
		<div class="my-16 flex flex-row px-4 text-2xl">
			<div class="grid grid-cols-2 gap-x-16 gap-y-8">
				<p>Name:</p>
				<p>{searchData?.Name ?? ''}</p>
				<p>Roll No.:</p>
				<p>{searchData?.RollNo ?? ''}</p>
				<p>Registration Status:</p>
				<p>{typeof searchData.Mess == 'number' ? 'Registered' : 'Unregistered'}</p>
				<p>Registered Mess:</p>
				<p>{getMessLabel(searchData.Mess)}</p>
				<p>Flagged</p>
				<p>No</p>
			</div>
			<div
				class="ml-24 flex w-max grow flex-col items-center justify-center border-l-2 border-l-custom-mid-grey"
			>
				<img src="" alt="profile pic" class="size-[12rem] rounded-full bg-white" />
			</div>
		</div>
		<div class="flex gap-x-4">
			<Modal bind:open={swapMessModalOpen} buttonText="Change Mess" class="ml-auto">
				<div class="flex flex-col gap-y-4 p-6">
					<h2 class="text-2xl font-semibold">Change Student Mess</h2>
					<p class="text-lg">
						Student: <span class="font-semibold">{searchData.Name}</span> ({searchData.RollNo})
					</p>
					<p class="text-lg">
						Current Mess: <span class="font-semibold">{getMessLabel(searchData.Mess)}</span>
					</p>
					<div class="mt-4 flex flex-col gap-y-2">
						<label for="mess-select" class="text-lg font-medium">Select New Mess:</label>
						<CustomSelect bind:value={selectedMess} items={messOptions} widthClass="w-full" />
					</div>
					<div class="mt-6 flex gap-x-4">
						<Button
							onclick={handleSwapMess}
							disabled={isSubmitting || !selectedMess}
							class="flex-1"
						>
							{isSubmitting ? 'Updating...' : 'Confirm Change'}
						</Button>
						<Button
							onclick={() => {
								swapMessModalOpen = false;
								// Reset to current mess
								if (typeof searchData.Mess === 'number') {
									selectedMess = searchData.Mess.toString();
								}
							}}
							disabled={isSubmitting}
							class="flex-1"
						>
							Cancel
						</Button>
					</div>
				</div>
			</Modal>

			<Modal bind:open={deregisterModalOpen} buttonText="Deregister">
				<div class="flex flex-col gap-y-4 p-6">
					<h2 class="text-2xl font-semibold">Deregister Student</h2>
					<p class="text-lg">
						Student: <span class="font-semibold">{searchData.Name}</span> ({searchData.RollNo})
					</p>
					<p class="text-lg">
						Current Mess: <span class="font-semibold">{getMessLabel(searchData.Mess)}</span>
					</p>

					<div class="mt-6 flex gap-x-4">
						<Button onclick={handleDeregister} disabled={isSubmitting} class="flex-1">
							{isSubmitting ? 'Deregistering...' : 'Confirm Deregister'}
						</Button>
						<Button
							onclick={() => {
								deregisterModalOpen = false;
							}}
							disabled={isSubmitting}
							class="flex-1"
						>
							Cancel
						</Button>
					</div>
				</div>
			</Modal>
		</div>
	{/if}
</section>
