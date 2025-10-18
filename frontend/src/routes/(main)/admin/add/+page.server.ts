import { PRIVATE_API_URL } from '$env/static/private';
import type { Actions } from './$types';

export const actions = {
	default: async ({ fetch, request }) => {
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
	}
} satisfies Actions;
