import { PUBLIC_API_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ parent, fetch }) => {
	const par = await parent();
	if (par.user['user_type'] != 0) {
		if (par.user['user_type'] == 1) {
			throw redirect(307, '/mess');
		} else {
			throw redirect(307, '/admin');
		}
	}
	return {
		regData: await (
			await fetch(PUBLIC_API_URL + '/students', {
				method: 'GET',
				credentials: 'include'
			})
		).json(),
		userMessData: await (
			await fetch(PUBLIC_API_URL + '/students/getMess', { method: 'GET', credentials: 'include' })
		).json()
	};
};
