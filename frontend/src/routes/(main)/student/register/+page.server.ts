import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { PRIVATE_API_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ parent }) => {
	const data = await parent();
	if (data.user.mess_id != 'No mess assigned') {
		throw redirect(307, '/');
	}
	return {
		messStats: await (
			await fetch(PRIVATE_API_URL + '/students/messStatsGrouped', {
				method: 'GET',
				credentials: 'include'
			})
		).json()
	};
};

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();
		if (data.get('veg') == 'true') {
			const res = await fetch(PRIVATE_API_URL + '/students/registerVegMess', {
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
				PRIVATE_API_URL + '/students/registerMess/' + parseInt(data.get('mess')?.toString()!),
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
