<script lang="ts">
	import { enhance } from '$app/forms';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Button from '$lib/components/common/Button.svelte';
	import Modal from '$lib/components/common/Modal.svelte';

	let { form } = $props();
	let fileInput = $state<HTMLInputElement>();
	let selectedFile = $state<File | null>(null);
	let isUploading = $state(false);
	let isDownloading = $state(false);
	let uploadModalOpen = $state(false);

	$effect(() => {
		if (form != undefined) {
			if (form.message != undefined) {
				if (form.success === false) {
					alert(form.message);
					if (form.errors && form.errors.length > 0) {
						console.error('CSV Upload Errors:', form.errors);
					}
				} else if (form.success === true) {
					let message = form.message;
					if (form.recordsAdded !== undefined) {
						message += `\n\nRecords Added: ${form.recordsAdded}`;
					}
					if (form.errors && form.errors.length > 0) {
						message += `\n\nErrors (${form.errors.length}):\n${form.errors.slice(0, 5).join('\n')}`;
						if (form.errors.length > 5) {
							message += `\n... and ${form.errors.length - 5} more errors`;
						}
					}
					alert(message);
					uploadModalOpen = false;
					setTimeout(() => {
						if (fileInput) {
							fileInput.value = '';
						}
						selectedFile = null;
					}, 0);
				} else {
					alert(form.message);
				}
			}
		}
	});

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		selectedFile = target.files?.[0] || null;
	}

	async function handleDownloadCSV() {
		isDownloading = true;
		try {
			const response = await fetch(PUBLIC_API_URL + '/office/students/download-csv', {
				method: 'GET',
				credentials: 'include'
			});

			if (response.ok) {
				const blob = await response.blob();
				const url = window.URL.createObjectURL(blob);
				const a = document.createElement('a');
				a.href = url;
				a.download = `students_${new Date().toISOString().split('T')[0]}.csv`;
				document.body.appendChild(a);
				a.click();
				window.URL.revokeObjectURL(url);
				document.body.removeChild(a);
			} else {
				alert('Failed to download CSV');
			}
		} catch (error) {
			console.error('Download error:', error);
			alert('Error downloading CSV');
		} finally {
			isDownloading = false;
		}
	}
</script>

<div class="flex flex-col dark:text-white">
	<div class="mb-8 flex items-center gap-x-4">
		<Modal bind:open={uploadModalOpen} buttonText="Upload Students CSV">
			<div class="flex min-w-[500px] flex-col gap-y-4 p-6">
				<h2 class="text-2xl font-semibold">Upload Students CSV</h2>
				<p class="text-sm text-gray-600 dark:text-gray-300">
					Expected CSV format: <span class="font-mono"
						>Name, Email, Phone, RollNo, Mess, Type, CanRegister</span
					>
				</p>
				<div class="rounded-md bg-gray-100 p-3 text-xs dark:bg-custom-dark-grey">
					<p class="mb-1">
						<strong>Mess Values:</strong> 1=Mess A-LDH, 2=Mess A-UDH, 3=Mess B-LDH, 4=Mess B-UDH, 5=Veg,
						0=None
					</p>
					<p class="mb-1"><strong>Type Values:</strong> 0=Student, 1=Staff, 2=Admin</p>
					<p><strong>CanRegister:</strong> true or false</p>
				</div>

				<form
					method="POST"
					action="?/uploadCSV"
					enctype="multipart/form-data"
					use:enhance={() => {
						isUploading = true;
						return async ({ update }) => {
							await update();
							isUploading = false;
						};
					}}
				>
					<div class="flex flex-col gap-y-4">
						<input
							bind:this={fileInput}
							type="file"
							name="file"
							accept=".csv"
							required
							onchange={handleFileSelect}
							class="w-full rounded-md bg-white px-3 py-2 text-sm file:mr-4 file:rounded-md file:border-0 file:bg-custom-light-orange file:px-4 file:py-2 file:text-sm file:font-medium file:text-white hover:file:cursor-pointer hover:file:brightness-95 dark:bg-custom-dark-grey"
						/>
						{#if selectedFile}
							<p class="text-sm text-gray-600 dark:text-gray-400">
								Selected: <span class="font-medium">{selectedFile.name}</span>
							</p>
						{/if}
						<div class="flex gap-x-4">
							<Button type="submit" disabled={isUploading || !selectedFile} class="flex-1">
								{isUploading ? 'Uploading...' : 'Upload CSV'}
							</Button>
							<Button
								type="button"
								onclick={() => {
									uploadModalOpen = false;
									setTimeout(() => {
										if (fileInput) {
											fileInput.value = '';
										}
										selectedFile = null;
									}, 0);
								}}
								disabled={isUploading}
								class="flex-1"
							>
								Cancel
							</Button>
						</div>
					</div>
				</form>
			</div>
		</Modal>

		<Button onclick={handleDownloadCSV} disabled={isDownloading}>
			{isDownloading ? 'Downloading...' : 'Download Students CSV'}
		</Button>
	</div>

	<hr class="my-6 border-gray-300 dark:border-gray-600" />

	<!-- Manual Add Student Form -->
	<h2 class="mb-6 text-2xl font-semibold">Add Single Student</h2>
	<form
		id="add-form"
		class="grid w-fit grid-rows-5 gap-y-6 text-lg *:grid *:grid-cols-2 *:grid-rows-subgrid **:[input]:w-[20rem] **:[input]:rounded-sm **:[input]:bg-white **:[input]:px-2.5 **:[input]:py-2 **:[input]:outline **:[input]:outline-black **:[input]:focus:ring-1 **:[input]:focus:ring-custom-light-orange **:[input]:focus:outline-none **:[input]:dark:bg-custom-mid-grey **:[input]:dark:outline-custom-dark-white/50"
		method="POST"
		action="?/addUser"
		use:enhance
	>
		<label>
			Name *
			<input required name="name" />
		</label>
		<label>
			Email *
			<input required name="email" type="email" />
		</label>
		<label>
			Roll Number *
			<input required name="roll" />
		</label>
		<label>
			Phone Number
			<input name="phone" />
		</label>
		<label>
			Designation *
			<select
				class="rounded-md bg-white px-2.5 py-1.5 outline-1 outline-black focus:ring-1 focus:ring-custom-light-orange focus:outline-none dark:bg-custom-mid-grey dark:outline-custom-dark-white/50"
				required
				name="type"
			>
				<option selected value={0}>Student</option>
				<option value={1}>Staff</option>
				<option value={2}>Admin</option>
			</select>
		</label>
		<label>
			Register to *
			<select
				class="rounded-md bg-white px-2.5 py-2 outline-1 outline-black focus:ring-1 focus:ring-custom-light-orange focus:outline-none dark:bg-custom-mid-grey dark:outline-custom-dark-white/50"
				required
				name="mess"
			>
				<option value={1}>Mess A - LDH</option>
				<option value={2}>Mess A - UDH</option>
				<option value={3}>Mess B - LDH</option>
				<option value={4}>Mess B - UDH</option>
				<option value={5}>Veg Mess</option>
				<option selected value={0}>None</option>
			</select>
		</label>
	</form>

	<Button
		type="submit"
		form="add-form"
		class="mt-12 self-center !bg-white !py-1.5 !text-lg dark:!bg-custom-mid-grey"
		>Confirm Registration</Button
	>
</div>
