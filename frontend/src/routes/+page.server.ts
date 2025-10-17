import { PUBLIC_API_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, parent }) => {
	const par = await parent();
	if (par.user['user_type'] == 0) {
		throw redirect(307, '/student');
	} else if (par.user['user_type'] == 1) {
		throw redirect(307, '/mess');
	} else if (par.user['user_type'] == 2) {
		throw redirect(307, '/admin');
	} else {
		await fetch(PUBLIC_API_URL + '/logout', { method: 'POST', credentials: 'include' });
		throw redirect(301, '/login');
	}
};
