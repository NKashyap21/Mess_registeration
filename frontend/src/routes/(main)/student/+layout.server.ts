import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { PRIVATE_API_URL } from '$env/static/private';

export const load: LayoutServerLoad = async ({ parent, fetch }) => {
	const par = await parent();
	if (par.user['user_type'] != 0) {
		if (par.user['user_type'] == 1) {
			throw redirect(307, '/mess');
		} else {
			throw redirect(307, '/admin');
		}
	}
	try {
		return {
			regData: await (
				await fetch(PRIVATE_API_URL + '/students/isRegistrationOpen', {
					method: 'GET',
					credentials: 'include'
				})
			).json(),
			userMessData: await (
				await fetch(PRIVATE_API_URL + '/students/getMess', {
					method: 'GET',
					credentials: 'include'
				})
			).json()
		};
	} catch (e) {
		console.error('Failed to fetch data under students route');
		console.error(e);
	}
};
