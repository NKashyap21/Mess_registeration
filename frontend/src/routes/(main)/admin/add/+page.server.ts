import { PRIVATE_API_URL } from '$env/static/private';
import type { Actions } from './$types';

export const actions = {
	addUser: async ({ fetch, request }) => {
		const formData = await request.formData();

		const jsonData = {
			name: formData.get('name')?.toString(),
			roll_no: formData.get('roll')?.toString(),
			phone: (formData.get('phone')?.toString() ?? '') == '' ? null : formData.get('phone'),
			email: formData.get('email')?.toString(),
			user_type: parseInt(formData.get('type')?.toString() ?? '0'),
			mess: parseInt(formData.get('mess')?.toString() ?? '0')
		};
		const res = await fetch(PRIVATE_API_URL + '/office/add-user', {
			method: 'POST',
			body: JSON.stringify(jsonData),
			credentials: 'include'
		});

		if (res.status == 200) {
			return {
				message: 'Successfully added the user.'
			};
		} else {
			return {
				message: `Failed to add user.\nReason: ${(await res.json()).error}`
			};
		}
	},

	uploadCSV: async ({ fetch, request }) => {
		const formData = await request.formData();
		const file = formData.get('file');

		if (!file || !(file instanceof File)) {
			return {
				success: false,
				message: 'No file uploaded'
			};
		}

		const uploadFormData = new FormData();
		uploadFormData.append('file', file);

		const res = await fetch(PRIVATE_API_URL + '/office/students/upload-csv', {
			method: 'POST',
			body: uploadFormData,
			credentials: 'include'
		});

		if (res.status === 200) {
			const data = await res.json();
			return {
				success: true,
				message: data.message,
				recordsAdded: data.records_added,
				errors: data.errors || []
			};
		} else {
			const error = await res.json();
			return {
				success: false,
				message: error.error || 'Failed to upload CSV'
			};
		}
	}
} satisfies Actions;
