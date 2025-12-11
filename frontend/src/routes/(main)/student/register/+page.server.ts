import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { PRIVATE_API_URL } from '$env/static/private';
import { resolve } from '$app/paths';

export const load: PageServerLoad = async ({ parent, fetch, url }) => {
	const data = await parent();
	if (
		data.userMessData?.data?.status == 'pending_sync' ||
		(data.userMessData?.data?.next_mess ?? 1) != 0
	) {
		throw redirect(307, resolve('/'));
	}
	return {
		messStats: await (
			await fetch(PRIVATE_API_URL + '/students/messStatsGrouped' + url.search, {
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
