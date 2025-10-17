import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { PRIVATE_API_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ fetch, parent }) => {
	const par = await parent();
	if (par.user['user_type'] == 0) {
		throw redirect(307, '/student');
	} else if (par.user['user_type'] == 1) {
		throw redirect(307, '/mess');
	} else if (par.user['user_type'] == 2) {
		throw redirect(307, '/admin');
	} else {
		await fetch(PRIVATE_API_URL + '/logout', { method: 'POST', credentials: 'include' });
		console.error('Failed page server load at /');
		throw redirect(307, '/login');
	}
};
