<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import Button from '$lib/components/common/Button.svelte';

	let rollNo = $state('');
	let searchData = $state<{ Name: string; RollNo: string; Mess: number | string }>({
		Name: '',
		Mess: '',
		RollNo: ''
	});
</script>

<section class="flex w-full flex-col">
	<div class="flex text-2xl">
		<input
			class="rounded-md border border-custom-mid-grey bg-white px-6 py-4"
			placeholder="Roll Number"
			bind:value={rollNo}
			onkeyup={(ev) => {
				if (ev.key == 'Enter') {
					fetch(PUBLIC_API_URL + '/office/students/' + rollNo, { credentials: 'include' })
						.then((res) => {
							if (res.status == 200) {
								res.json().then((data) => {
									searchData = data;
								});
							} else {
								alert('Error');
							}
						})
						.catch((e) => {
							console.error(e);
						});
				}
			}}
		/>
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
	<div class="my-16 flex flex-row px-4 text-2xl">
		<div class="grid grid-cols-2 gap-x-16 gap-y-8">
			<p>Name:</p>
			<p>{searchData?.Name ?? ''}</p>
			<p>Roll No.:</p>
			<p>{searchData?.RollNo ?? ''}</p>
			<p>Registration Status:</p>
			<p>{typeof searchData.Mess == 'number' ? 'Registered' : 'Unregistered'}</p>
			<p>Registered Mess:</p>
			<p>
				{Object.keys(searchData).includes('Mess')
					? searchData?.Mess == 1
						? 'Mess A - LDH'
						: searchData?.Mess == 2
							? 'Mess A - UDH'
							: searchData?.Mess == 3
								? 'Mess B - LDH'
								: searchData?.Mess == 4
									? 'Mess B - UDH'
									: 'Unknown'
					: 'Unregistered'}
			</p>
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
		<Button class="ml-auto">Swap Mess</Button>
		<Button>Deregister</Button>
	</div>
</section>
