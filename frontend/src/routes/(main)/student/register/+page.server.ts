import { PUBLIC_API_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ parent }) => {
	const data = await parent();
	if (data.user.mess_id != 'No mess assigned' || data.userMessData.status == 'pending_sync') {
		throw redirect(307, '/');
	}
};

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();
		if (data.get('veg') == 'true') {
			const res = await fetch(PUBLIC_API_URL + '/students/registerVegMess', {
				method: 'POST',
				credentials: 'include'
			});

			if (res.status != 200) {
				return {
					error: (await res.json()).error
				};
			}

			return {
				message: (await res.json()).message
			};
		} else {
			const res = await fetch(
				PUBLIC_API_URL + '/students/registerMess/' + parseInt(data.get('mess')?.toString()!),
				{
					method: 'POST',
					credentials: 'include'
				}
			);
			if (res.status != 200) {
				return {
					error: (await res.json()).error
				};
			}

			return {
				message: (await res.json()).message
			};
		}
	}
} satisfies Actions;
